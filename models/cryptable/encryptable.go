package encryptable

import (
	"fmt"

	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/utils/aes"
)

var aesKey = configurations.Configuration.Security.AESKEY

// Encryptable - public struct to implement sanitization by criptography
type Encryptable struct {
	Crypted bool `json:"Crypted"`
}

// SetAesKey to crypt
func SetAesKey(key string) {
	aesKey = key
}

// encrypt received value
func Encrypt(val string) (string, error) {
	encVal, err := aes.Encrypt(aesKey, val)
	if err != nil {
		return "", fmt.Errorf("encrypting: %w", err)
	}
	return encVal, nil
}

func Decrypt(val string) (string, error) {
	encVal, err := aes.Decrypt(aesKey, val)
	if err != nil {
		return "", fmt.Errorf("restoring: %v", err)
	}

	return encVal, nil
}

func (m *Encryptable) Encrypt(vals []*string) error {
	if m.Crypted {
		return nil
	}

	for i, val := range vals {
		encVal, err := aes.Encrypt(aesKey, *val)
		if err != nil {
			return fmt.Errorf("encrypting %d: %v", i, err)
		}
		*vals[i] = encVal
	}
	m.Crypted = true
	return nil
}

func (m *Encryptable) Restore(vals []*string) error {
	if !m.Crypted {
		return nil
	}

	for i, val := range vals {
		encVal, err := aes.Decrypt(aesKey, *val)
		if err != nil {
			return fmt.Errorf("restoring %v", err)
		}
		*vals[i] = encVal
	}
	m.Crypted = false
	return nil
}
