package enc_dec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

var encryptionKey = []byte{
	0x4f, 0x8a, 0x3d, 0x2c, 0x1b, 0x9e, 0x7d, 0x6f,
	0x5a, 0x0c, 0x8b, 0x4e, 0x3f, 0x2d, 0x1c, 0x9f,
	0x7e, 0x6d, 0x5b, 0x0d, 0x8c, 0x4f, 0x3e, 0x2e,
	0x1d, 0x9d, 0x7f, 0x6e, 0x5c, 0x0e, 0x8d, 0x4d,
}

func Encrypt(msg string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encryptedMsg string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func GenerateHash(msg string) string {
	hash := sha256.Sum256([]byte(msg))
	return base64.StdEncoding.EncodeToString(hash[:])
}
