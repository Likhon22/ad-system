package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", errors.New("failed to generate salt")
	}
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	saltHash64 := base64.StdEncoding.EncodeToString(hash)
	encodedString := fmt.Sprintf("%s.%s", saltBase64, saltHash64)
	return encodedString, nil
}

func ComparePassword(inputPass, storedHash string) (bool, error) {
	parts := strings.SplitN(storedHash, ".", 2)
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}
	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, errors.New("failed to decode salt")
	}
	storedHashBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, errors.New("failed to decode hash")
	}
	inputHash := argon2.IDKey([]byte(inputPass), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	match := subtle.ConstantTimeCompare(storedHashBytes, inputHash)
	return match == 1, nil

}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.StdEncoding.EncodeToString(hash[:])
}
