package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

type Encryptor interface {
	Encrypt(string) string
	Decrypt(string) (string, error)
}

type encryptor struct {
	salt     []byte
	password []byte
}

var enc = new(encryptor)

func Init(salt string, password string) {
	bit16Salt := md5.Sum([]byte(salt)) //nolint:gosec
	passwordBytes := sha256.Sum256([]byte(password))
	enc = &encryptor{
		salt:     bit16Salt[:],
		password: passwordBytes[:],
	}
}

func GetEncryptor() Encryptor {
	return enc
}

func (e *encryptor) Encrypt(text string) string {
	plainText := []byte(text)
	cipherText := make([]byte, len(plainText))
	block, err := aes.NewCipher(e.password)
	if err != nil {
		log.Println("Cannot initialize encryptor")
		return ""
	}
	encrypter := cipher.NewCFBEncrypter(block, e.salt)
	encrypter.XORKeyStream(cipherText, plainText)
	return encode(cipherText)
}

func (e *encryptor) Decrypt(text string) (string, error) {
	cipherText, err := decode(text)
	if err != nil {
		return "", err
	}
	plainText := make([]byte, len(cipherText))
	block, err := aes.NewCipher(e.password)
	if err != nil {
		log.Println("Cannot initialize encryptor")
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, e.salt)
	decrypter.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func encode(b []byte) string {
	return hex.EncodeToString(b)
}

func decode(s string) ([]byte, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}
