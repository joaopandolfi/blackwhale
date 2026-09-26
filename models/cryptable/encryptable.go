package encryptable

import (
	"fmt"

	"github.com/joaopandolfi/blackwhale/v2/configurations"
	"github.com/joaopandolfi/blackwhale/v2/utils/aes"
)

var aesKeyOverride string

// Encryptable - public struct to implement sanitization by criptography
type Encryptable struct {
	Crypted bool
}

// SetAesKey to crypt
func SetAesKey(key string) {
	aesKeyOverride = key
}

func key() string {
	if aesKeyOverride != "" {
		return aesKeyOverride
	}
	return configurations.Configuration.Security.AESKEY
}

// encrypt received value
func Encrypt(val string) (string, error) {
	encVal, err := aes.Encrypt(key(), val)
	if err != nil {
		return "", fmt.Errorf("encrypting: %w", err)
	}
	return encVal, nil
}

func Decrypt(val string) (string, error) {
	encVal, err := aes.Decrypt(key(), val)
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
		encVal, err := aes.Encrypt(key(), *val)
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
		encVal, err := aes.Decrypt(key(), *val)
		if err != nil {
			return fmt.Errorf("restoring %v", err)
		}
		*vals[i] = encVal
	}
	m.Crypted = false
	return nil
}
