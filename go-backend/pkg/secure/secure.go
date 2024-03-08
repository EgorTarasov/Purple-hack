package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"purple/internal/shared"
)

func Encrypt(toEncrypt, keyString string) (string, error) {
	key, _ := hex.DecodeString(keyString)
	plain := []byte(toEncrypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plain))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plain)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(toDecrypt, keyString string) (string, error) {
	key, _ := hex.DecodeString(keyString)
	cipherText, _ := base64.URLEncoding.DecodeString(toDecrypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", shared.ErrCipherTooShort
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}
