package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

const encryptedPrefix = "enc:"

func EncryptString(secret, value string) (string, error) {
	if value == "" {
		return "", nil
	}

	block, err := aes.NewCipher(key(secret))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	sealed := gcm.Seal(nonce, nonce, []byte(value), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(sealed), nil
}

func DecryptString(secret, value string) (string, error) {
	if value == "" || !strings.HasPrefix(value, encryptedPrefix) {
		return value, nil
	}

	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, encryptedPrefix))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key(secret))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(payload) < gcm.NonceSize() {
		return "", errors.New("encrypted payload is too short")
	}

	nonce := payload[:gcm.NonceSize()]
	ciphertext := payload[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func key(secret string) []byte {
	sum := sha256.Sum256([]byte(secret))
	return sum[:]
}
