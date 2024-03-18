package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"os"
    "strconv"
)

func HashPassword(p string) (string, error) {
    bcryptSalt := os.Getenv("BCRYPT_SALT")
    cost, err := strconv.Atoi(bcryptSalt)
    if err != nil {
        // Handle error if BCRYPT_SALT is not a valid integer
        return "", err
    }

    password := []byte(p)
    hash, err := bcrypt.GenerateFromPassword(password, cost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func ComparePassword(h, p []byte) bool {
	hash, pass := []byte(h), []byte(p)
	err := bcrypt.CompareHashAndPassword(hash, pass)
	return err == nil
}
