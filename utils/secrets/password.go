package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/private-square/bkst-users-api/utils/errors"
	"io"
	mr "math/rand"
	"time"
	"unicode"
)

// GetRandomPassword generates a random string of upper + lower case alphabets and digits
// which is 23 bits long and returns the string
func GetRandomPassword() string {
	mr.Seed(time.Now().UnixNano())
	digits := "0123456789"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + digits
	length := 23
	buf := make([]byte, length)
	buf[0] = digits[mr.Intn(len(digits))]
	for i := 1; i < length; i++ {
		buf[i] = all[mr.Intn(len(all))]
	}
	mr.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

func VerifyPassword(password string) error {
	var (
		numOfLetters                  = 0
		number, upper, lower, special bool
	)
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
			numOfLetters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			numOfLetters++
		case unicode.IsUpper(c):
			upper = true
			numOfLetters++
		case unicode.IsLower(c) || c == ' ':
			lower = true
			numOfLetters++
		}
	}
	if numOfLetters > 8 && number && upper && lower && special {
		return nil
	} else {
		return errors.InvalidPasswordError
	}
}

func createSHA256Hash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func EncryptPassword(data, passphrase string) (string, error) {
	block, _ := aes.NewCipher(createSHA256Hash(passphrase))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.PasswordEncryptionError{Err: err}
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.PasswordEncryptionError{Err: err}
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptPassword(data, passphrase string) (string, error) {
	bData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.PasswordDecryptionError{Err: err}
	}
	key := createSHA256Hash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.PasswordDecryptionError{Err: err}
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.PasswordDecryptionError{Err: err}
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := bData[:nonceSize], bData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.PasswordDecryptionError{Err: err}
	}
	return string(plaintext), nil
}
