package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"hack/internal/shared"
	"io"
	"sync"
)

var mu sync.Mutex

func GetPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

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

func GeneratePassResetCode() (string, error) {
	codeRaw := make([]byte, 3)
	if _, err := rand.Read(codeRaw); err != nil {
		return "", err
	}
	code := hex.EncodeToString(codeRaw)
	return code, nil
}

func GenerateHash(src string) string {
	mu.Lock()
	hash := sha256.Sum256([]byte(src))
	mu.Unlock()
	return base64.StdEncoding.EncodeToString(hash[:])
}
