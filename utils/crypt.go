package utils

import (
	"github.com/joaopandolfi/blackwhale/v2/configurations"
	jwt "github.com/joaopandolfi/blackwhale/v2/remotes/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Token is an alias of the canonical jwt.Token
type Token = jwt.Token

// HashPassword - Make password hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(setSecretOnPass(password)), configurations.Configuration.Security.BCryptCost)
	return string(bytes), err
}

// CheckPasswordHash - Chek if password and hash is correspondent
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(setSecretOnPass(password)))
	return err == nil
}

func setSecretOnPass(password string) string {
	return configurations.Configuration.BCryptSecret + "!" + password
}

// CheckJwtToken - Check sended token using the configured jwt secret
func CheckJwtToken(tokenString string) (Token, error) {
	return jwt.CheckJwtToken(tokenString, configurations.Configuration.Security.JWTSecret)
}

// NewJwtToken - Crete token with expiration time
func NewJwtToken(t Token, expMinutes int) (string, error) {
	t.Authorized = true
	return NewJwtTokenV2(t, expMinutes)
}

// NewJwtTokenV2 - Create a token with expiration time using the configured jwt secret
func NewJwtTokenV2(t Token, expMinutes int) (string, error) {
	return jwt.NewJwtToken(t, expMinutes, configurations.Configuration.Security.JWTSecret)
}
