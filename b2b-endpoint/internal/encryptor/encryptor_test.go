package encryptor

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEncryptor_EncryptDecrypt(t *testing.T) {
	key := "12345678901234567890123456789012" // 32-байтный ключ для AES-256
	enc := NewEncryptor(key)

	plainText := uuid.New().String()

	// Шифрование
	encryptedText, err := enc.Encrypt(plainText)
	require.NoError(t, err, "Encryption failed")
	require.NotEmpty(t, encryptedText, "Encrypted text should not be empty")

	// Расшифрование
	decryptedText, err := enc.Decrypt(encryptedText)
	require.NoError(t, err, "Decryption failed")
	require.Equal(t, plainText, decryptedText, "Decrypted text does not match original")
}

func TestEncryptor_DecryptInvalidData(t *testing.T) {
	key := "12345678901234567890123456789012"
	enc := NewEncryptor(key)

	// Попытка расшифровать некорректные данные
	invalidEncryptedText := "invalid_cipher_text"
	_, err := enc.Decrypt(invalidEncryptedText)
	require.Error(t, err, "Decryption of invalid data should fail")
}
