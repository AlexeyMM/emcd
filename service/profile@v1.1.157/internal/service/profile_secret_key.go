package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/pkg/aes"
)

const (
	keyLength              = 32
	opGetOrCreateSecretKey = "server.profile.GetOrCreateSecretKey"
	opCheckSignature       = "server.profile.CheckSignature"
)

type GetOrCreateSecretKeyReq struct {
	UserUUID, ParentUUID uuid.UUID
	IP                   string
}

func (p *profile) GetOrCreateSecretKey(ctx context.Context, req GetOrCreateSecretKeyReq) (string, error) {

	user, err := p.profileRepo.GetByUserID(ctx, req.UserUUID)
	if err != nil {
		log.Error(ctx, "%s: get old user id: %s", opGetOrCreateSecretKey, err.Error())
		return "", fmt.Errorf("%w: get old user id by uuid", ErrNotFound)
	}

	if user.User.SecretKey != "" {
		secretKey, err := decrypt(user.User.SecretKey)
		if err != nil {
			log.Error(ctx, "%s: decrypt secret key: %s", opGetOrCreateSecretKey, err.Error())
			return "", fmt.Errorf("%w: decrypt secret key", ErrInternal)
		}
		return secretKey, nil
	}

	newSecretKey, err := generateSecretKey()
	if err != nil {
		log.Error(ctx, "%s: generate secret key: %s", opGetOrCreateSecretKey, err.Error())
		return "", err
	}

	err = p.profileRepo.SetSecretKeyForUserID(ctx, encrypt(newSecretKey), req.UserUUID)
	if err != nil {
		log.Error(ctx, "%s: set secret key for user: %s", opGetOrCreateSecretKey, err.Error())
		return "", err
	}

	return newSecretKey, nil
}

type CheckSignatureReq struct {
	UserUUID, ParentUUID uuid.UUID
	Signature            string
	Nonce                int64
	IP                   string
}

func (p *profile) CheckSignature(ctx context.Context, req CheckSignatureReq) (bool, error) {
	apiKey, err := p.GetAPIKey(ctx, GetAPIKeyReq{
		UserUUID:   req.UserUUID,
		ParentUUID: req.ParentUUID,
		IP:         req.IP,
	})
	if err != nil {
		log.Error(ctx, "%s: get api key: %s", opCheckSignature, err.Error())
		return false, err
	}
	u, err := p.profileRepo.GetByUserID(ctx, req.UserUUID)
	if err != nil {
		log.Error(ctx, "%s: get user: %s", opCheckSignature, err.Error())
		return false, err
	}
	secretKey, err := decrypt(u.User.SecretKey)
	if err != nil {
		log.Error(ctx, "%s: decrypt secret key: %s", opCheckSignature, err.Error())
		return false, err
	}

	check := verifySignature(apiKey,
		strconv.Itoa(int(req.Nonce)),
		req.UserUUID.String(),
		secretKey,
		req.Signature)
	if !check {
		log.Error(ctx, "check signature is false: nonce %d, user %s",
			req.Nonce, req.UserUUID.String())
		return false, nil
	}

	checkNonce, err := p.nonceStore.CheckAndUpdateNonce(ctx, req.UserUUID.String(), req.Nonce)
	if err != nil {
		log.Error(ctx, "%s: check and update nonce: %s", opCheckSignature, err.Error())
		return false, err
	}
	if !checkNonce {
		log.Error(ctx, "%s: check and update nonce: nonce is invalid", opCheckSignature)
		return false, nil
	}
	return check, nil
}

func encrypt(key string) string {
	return aes.GetEncryptor().Encrypt(key)
}

func decrypt(key string) (string, error) {
	return aes.GetEncryptor().Decrypt(key)
}

func generateSecretKey() (string, error) {
	key := make([]byte, keyLength)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}

func verifySignature(apiKey, nonce, userID, secretKey, providedSignature string) bool {
	message := userID + apiKey + nonce
	key := []byte(secretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	computedSignature := hex.EncodeToString(h.Sum(nil))

	computedSignature = strings.ToUpper(computedSignature)

	return computedSignature == providedSignature
}
