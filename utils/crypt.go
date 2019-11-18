package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"octopus.com/mi/configurations"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(setSecretOnPass(password)), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(setSecretOnPass(password)))
	return err == nil
}

func setSecretOnPass(password string) string {
	return fmt.Sprintf("%s!%s", configurations.Configuration.BCryptSecret, password)
}
