package configurations

import (
	"encoding/json"
	"fmt"

	"github.com/joaopandolfi/blackwhale/remotes/request"
	"golang.org/x/xerrors"
)

type VaultPayload struct {
	Success bool              `json:"success"`
	Data    map[string]string `json:"data"`
	Message string            `json:"message"`
}

func LoadVault(host, token, privKey string) (map[string]string, error) {
	var result VaultPayload
	server := fmt.Sprintf("%s/vault/recover", host)

	b, err := request.GetWithHeader(server, map[string]string{
		"key":   privKey,
		"token": token,
	})
	if err != nil {
		return nil, xerrors.Errorf("loading vault from server (%s): %w", server, err)
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, xerrors.Errorf("unmarshaling payload response: %w", err)
	}

	if !result.Success {
		return nil, xerrors.Errorf("received error on load vault: %s", result.Message)
	}

	return result.Data, nil
}
