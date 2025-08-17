package auth

import (
	"encoding/hex"
	"fmt"
	"strings"
	"crypto/rand"
	"crypto/subtle"
	"golang.org/x/crypto/argon2"
)

const (
	time_   = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
)

func HashPassword(password string) string {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	hash := argon2.IDKey([]byte(password), salt, time_, memory, threads, keyLen)
	return fmt.Sprintf("%x.%x", salt, hash)
}

func CheckPassword(password, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hash format")
	}
	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	storedHash, err := hex.DecodeString(parts[1])
	if err != nil {
		return false, err
	}
	newHash := argon2.IDKey([]byte(password), salt, time_, memory, threads, keyLen)
	return subtle.ConstantTimeCompare(storedHash, newHash) == 1, nil
}
