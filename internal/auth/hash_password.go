package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

var (
	ErrCheckingWithArgon = "error while chceking the password hash"
	ErrHashingWithArgon  = "error while hashing password"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("%v: %v", ErrHashingWithArgon, err)
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return match, fmt.Errorf("%v: %v", ErrCheckingWithArgon, err)
	}

	return match, nil
}
