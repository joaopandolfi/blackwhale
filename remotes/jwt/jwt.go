package jwt

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// NewJwtToken - Crete token with expiration time
func NewJwtToken(t Token, expMinutes int, secret string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = t.Authorized
	atClaims["id"] = t.ID
	atClaims["institution"] = t.Institution
	atClaims["permissions"] = t.Permission
	atClaims["broker"] = t.Broker
	atClaims["exp"] = time.Now().Add(time.Minute * time.Duration(expMinutes)).Unix()

	return newJwt(atClaims, secret)
}

// NewJwtGeneric - Create an jwt using custom data
func NewJwtGeneric(data map[string]interface{}, expMinutes int, secret string) (string, error) {
	atClaims := jwt.MapClaims{}
	for k, v := range data {
		atClaims[k] = v
	}
	atClaims["exp"] = time.Now().Add(time.Minute * time.Duration(expMinutes)).Unix()

	return newJwt(atClaims, secret)
}

func newJwt(atClaims jwt.MapClaims, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// CheckJwtToken - Check sended token
func CheckJwtToken(tokenString, secret string) (Token, error) {

	token, err := CheckJwtGenericToken(tokenString, secret)
	if err != nil {
		return Token{Authorized: false}, err
	}

	return Token{
		Authorized:  true,
		ID:          token["id"].(string),
		Institution: token["institution"].(string),
		Permission:  token["permissions"].(string),
		Broker:      token["broker"].(map[string]interface{}),
	}, nil
}

func CheckJwtGenericToken(tokenString, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing jwt: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid Token")
	}

	exps := claims["exp"].(float64)
	if int64(exps) < time.Now().Unix() {
		return nil, fmt.Errorf("expired token")
	}

	return claims, nil
}
