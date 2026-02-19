package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", errors.New("failed to generate salt")
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	saltHash64 := base64.StdEncoding.EncodeToString(hash)
	encodedString := fmt.Sprintf("%s.%s", saltBase64, saltHash64)
	return encodedString, nil
}
