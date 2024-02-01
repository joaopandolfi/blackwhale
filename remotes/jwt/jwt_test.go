package jwt

import "testing"

func TestJWT(t *testing.T) {
	secret := "secret-super-secret-with-32-char"
	baseToken := Token{
		ID:          "id",
		Permission:  "permission",
		Institution: "institution",
		Broker: map[string]interface{}{
			"test": "test",
		},
	}

	jwt, err := NewJwtToken(baseToken, 5, secret)
	if err != nil {
		t.Errorf("creating jwt: %v", err)
		return
	}

	token, err := CheckJwtToken(jwt, secret)
	if err != nil {
		t.Errorf("checking jwt token: %v", err)
		return
	}

	if token.ID != baseToken.ID {
		t.Errorf("base token token id expected id got: %v", token.ID)
		return
	}
}
