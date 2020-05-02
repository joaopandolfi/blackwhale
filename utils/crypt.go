package utils

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joaopandolfi/blackwhale/configurations"
	"golang.org/x/crypto/bcrypt"
)

// Token -
type Token struct {
	ID          string `json:'id'`
	Permission  string `json:'permission'`
	Institution string `json:'institution'`
	Authorized  bool   `json:'authorized'`
}

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
	return fmt.Sprintf("%s!%s", configurations.Configuration.BCryptSecret, password)
}

// CheckJwtToken - Check sended token
func CheckJwtToken(tokenString string) (Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(configurations.Configuration.Security.JWTSecret), nil
	})
	if err != nil {
		return Token{Authorized: false}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return Token{Authorized: false}, fmt.Errorf("invalid Token")
	}

	exps := claims["exp"].(string)
	exp, _ := strconv.ParseInt(exps, 10, 64)
	if exp < time.Now().Unix() {
		return Token{Authorized: false}, fmt.Errorf("expired token")
	}

	return Token{
		Authorized:  true,
		ID:          claims["id"].(string),
		Institution: claims["institution"].(string),
		Permission:  claims["permission"].(string),
	}, nil
}

// NewJwtToken - Crete token with expiration time
func NewJwtToken(t Token, expMinutes int) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = t.ID
	atClaims["institution"] = t.Institution
	atClaims["permission"] = t.Permission
	atClaims["exp"] = fmt.Sprintf("%v", time.Now().Add(time.Minute*time.Duration(expMinutes)).Unix())
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(configurations.Configuration.Security.JWTSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}
