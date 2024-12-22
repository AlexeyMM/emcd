package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/model"
	"code.emcdtech.com/emcd/service/profile/internal/pkg/aes"
)

const (
	userLogChangeTypeCreateApiKey = "newkey"
	userLogChangeTypeDeleteApiKey = "revokeapi"
	userLogGetAPIKey              = "getapikey"
	opCreateAPIKey                = "server.profile.CreateAPIKey"
	opDeleteAPIKey                = "server.profile.DeleteAPIKey"
	opGetAPIKey                   = "server.profile.opGetAPIKey"
)

type CreateNewAPIKeyReq struct {
	UserUUID, ParentUUID uuid.UUID
	IP                   string
}

func (p *profile) CreateNewAPIKey(ctx context.Context, req CreateNewAPIKeyReq) (string, error) {
	var apiKey string
	var err error
	var user, parent *model.User

	user, err = p.oldUsers.GetUserByUUID(ctx, req.UserUUID)
	if err != nil {
		log.Error(ctx, "%s: get old user id: %s", opCreateAPIKey, err.Error())
		return "", fmt.Errorf("%w: get old user id by uuid", ErrNotFound)
	}

	if req.ParentUUID != uuid.Nil {
		parent, err = p.oldUsers.GetUserByUUID(ctx, req.ParentUUID)
		if err != nil {
			log.Error(ctx, "%s: get old parent id: %s", opCreateAPIKey, err.Error())
			return "", fmt.Errorf("%w: get old parent id by uuid", ErrNotFound)
		}

		apiKey, err = p.createAPIKeyByUserUUIDAndParentUUID(ctx, user.OldID, parent.OldID)
	} else {
		apiKey, err = p.createAPIKeyByUserUUID(ctx, user.OldID)
	}

	if err != nil {
		return "", err
	}

	apiKey, err = aes.GetEncryptor().Decrypt(apiKey)
	if err != nil {
		log.Error(ctx, "%s: decrypt api key for user: %d, err: %s", opCreateAPIKey, user.OldID, err.Error())
		return apiKey, fmt.Errorf("%w: decrypt api key", ErrInternal)
	}

	err = p.userLogs.CreateWithoutToken(ctx, &model.UserLog{
		UserID:     user.OldID,
		ChangeType: userLogChangeTypeCreateApiKey,
		IP:         req.IP,
		OldValue:   "",
		Value:      apiKey,
		Active:     true,
	})
	if err != nil {
		log.Error(ctx, "%s: write to user log: %s", opCreateAPIKey, err.Error())
	}
	return apiKey, err
}

func (p *profile) createAPIKeyByUserUUIDAndParentUUID(ctx context.Context, userID, parentID int32) (string, error) {
	log.Info(ctx, "%s: start create api key for user: %d, parent: %d", opCreateAPIKey, userID, parentID)
	apiKeyExist, _, err := p.getAPIKeyByUserUUIDAndParentUUID(ctx, userID, parentID)
	if err != nil {
		log.Error(ctx, "%s: get api key for user: %d and parent: %d, err: %s",
			opCreateAPIKey, userID, parentID, err.Error())
		return "", fmt.Errorf("%w: get api key by user and parent: %s", ErrInternal, err.Error())
	}

	if apiKeyExist {
		log.Error(ctx, "%s: api key for user: %d already exist", opCreateAPIKey, userID)
		return "", fmt.Errorf("%w: api key already exist", ErrAlreadyExist)
	}

	apiRetKey := uuid.New().String()
	apiKey := aes.GetEncryptor().Encrypt(apiRetKey)

	err = p.setAPIKeyForUserUUIDAndParentUUID(ctx, apiKey, userID, parentID)
	if err != nil {
		log.Error(ctx,
			"%s: set api key for pair user_id: %d and parent_id: %d: %s",
			opCreateAPIKey, userID, parentID, err.Error())
		return "", fmt.Errorf("%w: set api key by pair user_id and parent_id: %s", ErrInternal, err.Error())
	}
	log.Info(ctx, "%s: finish create api key for user: %d, parent: %d", opCreateAPIKey, userID, parentID)

	return apiKey, nil
}

func (p *profile) createAPIKeyByUserUUID(ctx context.Context, userID int32) (string, error) {
	log.Info(ctx, "%s: start create api key for user: %d", opCreateAPIKey, userID)
	apiKeyExist, _, err := p.getAPIKeyByUserUUID(ctx, userID)
	if err != nil {
		log.Error(ctx, "%s: get api key for user: %d, err: %s", opCreateAPIKey, err.Error())
		return "", fmt.Errorf("%w: get api key by user id: %s", ErrInternal, err.Error())
	}
	if apiKeyExist {
		log.Error(ctx, "%s: api key for user: %d already exist", opCreateAPIKey, userID)
		return "", fmt.Errorf("%w: api key already exist", ErrAlreadyExist)
	}

	apiRetKey := uuid.New().String()
	apiKey := aes.GetEncryptor().Encrypt(apiRetKey)

	err = p.setAPIKeyForUserUUID(ctx, apiKey, userID)
	if err != nil {
		log.Error(ctx, "%s: set api key for user_id: %d: %s", opCreateAPIKey, userID, err.Error())
		return "", fmt.Errorf("%w: set api key by user_id: %s", ErrInternal, err.Error())
	}
	log.Info(ctx, "%s: finish create api key for user: %d", opCreateAPIKey, userID)
	return apiKey, nil
}

type DeleteAPIKeyReq struct {
	UserUUID, ParentUUID uuid.UUID
	IP                   string
}

func (p *profile) DeleteAPIKey(ctx context.Context, req DeleteAPIKeyReq) error {
	var apiKey string
	var err error
	var user, parent *model.User

	user, err = p.oldUsers.GetUserByUUID(ctx, req.UserUUID)
	if err != nil {
		log.Error(ctx, "%s: get old user id: %s", opDeleteAPIKey, err.Error())
		return fmt.Errorf("%w: get old user id by uuid", ErrNotFound)
	}

	if req.ParentUUID != uuid.Nil {
		parent, err = p.oldUsers.GetUserByUUID(ctx, req.ParentUUID)
		if err != nil {
			log.Error(ctx, "%s: get old parent id: %s", opDeleteAPIKey, err.Error())
			return fmt.Errorf("%w: get old parent id by uuid", ErrNotFound)
		}
	}

	var oldParentID int32
	if parent != nil {
		oldParentID = parent.OldID
	}

	apiKey, err = p.getOldAPIKey(ctx, user.OldID, oldParentID)
	if err != nil {
		return err
	}

	if req.ParentUUID != uuid.Nil {
		err = p.deleteAPIKeyByUserUUIDAndParentUUID(ctx, user.OldID, oldParentID)
	} else {
		err = p.deleteAPIKeyByUserUUID(ctx, user.OldID)
	}

	if err != nil {
		return err
	}

	err = p.userLogs.CreateWithoutToken(ctx, &model.UserLog{
		UserID:     user.OldID,
		ChangeType: userLogChangeTypeDeleteApiKey,
		IP:         req.IP,
		OldValue:   apiKey,
		Value:      "",
		Active:     true,
	})
	if err != nil {
		log.Error(ctx, "%s: write to user log: %s", opDeleteAPIKey, err.Error())
	}

	return nil
}

type GetAPIKeyReq struct {
	UserUUID, ParentUUID uuid.UUID
	IP                   string
}

func (p *profile) GetAPIKey(ctx context.Context, req GetAPIKeyReq) (string, error) {
	var apiKey string
	var err error
	var apiKeyExist bool
	var user, parent *model.User

	user, err = p.oldUsers.GetUserByUUID(ctx, req.UserUUID)
	if err != nil {
		log.Error(ctx, "%s: get old user id: %s", opGetAPIKey, err.Error())
		return "", fmt.Errorf("%w: get old user id by uuid", ErrNotFound)
	}

	if req.ParentUUID != uuid.Nil {
		parent, err = p.oldUsers.GetUserByUUID(ctx, req.ParentUUID)
		if err != nil {
			log.Error(ctx, "%s: get old parent id: %s", opGetAPIKey, err.Error())
			return "", fmt.Errorf("%w: get old parent id by uuid", ErrNotFound)
		}

		apiKeyExist, apiKey, err = p.getAPIKeyByUserUUIDAndParentUUID(ctx, user.OldID, parent.OldID)
	} else {
		apiKeyExist, apiKey, err = p.getAPIKeyByUserUUID(ctx, user.OldID)
	}

	if err != nil {
		log.Error(ctx, "%s: get api key for user: %d, err: %s", opGetAPIKey, user.OldID, err.Error())
		return apiKey, fmt.Errorf("%w: get api key by user: %s", ErrInternal, err.Error())
	}

	if !apiKeyExist {
		log.Error(ctx, "%s: no api key for user: %d", opGetAPIKey, user.OldID)
		return apiKey, fmt.Errorf("%w: no api key", ErrNotFound)
	}

	apiKey, err = aes.GetEncryptor().Decrypt(apiKey)
	if err != nil {
		log.Error(ctx, "%s: decrypt api key for user: %d, err: %s", opGetAPIKey, user.OldID, err.Error())
		return apiKey, fmt.Errorf("%w: decrypt api key", ErrInternal)
	}

	err = p.userLogs.CreateWithoutToken(ctx, &model.UserLog{
		UserID:     user.OldID,
		ChangeType: userLogGetAPIKey,
		IP:         req.IP,
		OldValue:   "",
		Value:      "",
		Active:     true,
	})
	if err != nil {
		log.Error(ctx, "%s: write to user log: %s", opGetAPIKey, err.Error())
	}

	return apiKey, nil
}

func (p *profile) getOldAPIKey(ctx context.Context, userID, parentID int32) (string, error) {
	var oldAPIKey string
	var err error
	var apiKeyExist bool

	if parentID > 0 {
		apiKeyExist, oldAPIKey, err = p.getAPIKeyByUserUUIDAndParentUUID(ctx, userID, parentID)
	} else {
		apiKeyExist, oldAPIKey, err = p.getAPIKeyByUserUUID(ctx, userID)
	}

	if err != nil {
		log.Error(ctx, "%s: get api key for user: %d and parent: %d, err: %s",
			opDeleteAPIKey, userID, parentID, err.Error())
		return oldAPIKey, fmt.Errorf("%w: get api key by user and parent: %s", ErrInternal, err.Error())
	}

	if !apiKeyExist {
		log.Error(ctx, "%s: no api key for user: %d", opDeleteAPIKey, userID)
		return oldAPIKey, fmt.Errorf("%w: no api key", ErrNotFound)
	}

	return oldAPIKey, nil
}

func (p *profile) deleteAPIKeyByUserUUIDAndParentUUID(ctx context.Context, userID, parentID int32) error {
	err := p.setAPIKeyForUserUUIDAndParentUUID(ctx, "", userID, parentID)
	if err != nil {
		log.Error(ctx,
			"%s: delete api key for pair user_id: %s and parent_id: %d: %d",
			opDeleteAPIKey, userID, parentID, err.Error())
		return fmt.Errorf("%w: delete api key by pair user_id and parent_id: %s", ErrInternal, err.Error())
	}
	return nil
}

func (p *profile) deleteAPIKeyByUserUUID(ctx context.Context, userID int32) error {
	err := p.setAPIKeyForUserUUID(ctx, "", userID)
	if err != nil {
		log.Error(ctx, "%s: delete api key for user_id: %d: %s", opDeleteAPIKey, userID, err.Error())
		return fmt.Errorf("%w: delete api key by user_id: %s", ErrInternal, err.Error())
	}
	return nil
}

func (p *profile) getAPIKeyByUserUUIDAndParentUUID(ctx context.Context, userID, parentID int32) (bool, string, error) {
	exist, existedAPIKey, err := p.oldUsers.GetAPIKeyByUserIDAndParentID(ctx, userID, parentID)
	if err != nil {
		return false, "", fmt.Errorf("%w: get api key by user id and parent id", err)
	}
	return exist, existedAPIKey, nil
}

func (p *profile) setAPIKeyForUserUUIDAndParentUUID(ctx context.Context, apiKey string, userID, parentID int32) error {
	err := p.oldUsers.SetAPIKeyForUserIDAndParentID(ctx, apiKey, userID, parentID)
	if err != nil {
		return fmt.Errorf("%w: set api key for pair user_id and parent_id", err)
	}
	return nil
}

func (p *profile) getAPIKeyByUserUUID(ctx context.Context, userID int32) (bool, string, error) {
	var exist bool
	u, err := p.oldUsers.GetUserByOldID(ctx, int(userID))
	if err != nil {
		return false, "", fmt.Errorf("%w: get api key by user id", err)
	}
	if u.ApiKey != "" {
		exist = true
	}
	return exist, u.ApiKey, nil
}

func (p *profile) setAPIKeyForUserUUID(ctx context.Context, apiKey string, userID int32) error {
	err := p.oldUsers.SetAPIKeyForUserID(ctx, apiKey, userID)
	if err != nil {
		return fmt.Errorf("%w: set api key for user_id", err)
	}
	return nil
}
