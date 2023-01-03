package sanitizable

import (
	"fmt"

	config "github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/utils/aes"
)

// Sanitilizable - public struct to implement sanitization by criptography
type Sanitilizable struct {
	Sanitized bool `json:"sanitized" `
}

func (m *Sanitilizable) Sanitize(vals map[string]*string) error {
	if m.Sanitized {
		return nil
	}

	for key, val := range vals {
		encVal, err := aes.Encrypt(config.Configuration.Security.AESKEY, *val)
		if err != nil {
			return fmt.Errorf("encrypting %s: %v", key, err)
		}
		*vals[key] = encVal
	}
	m.Sanitized = true
	return nil
}

func (m *Sanitilizable) Restore(vals map[string]*string) error {
	if !m.Sanitized {
		return nil
	}

	for key, val := range vals {
		encVal, err := aes.Decrypt(config.Configuration.Security.AESKEY, *val)
		if err != nil {
			return fmt.Errorf("restoring %s: %v", key, err)
		}
		*vals[key] = encVal
	}
	m.Sanitized = false
	return nil
}
