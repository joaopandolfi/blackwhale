//go:build integration

package hasura_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/joaopandolfi/blackwhale/v2/remotes/hasura"
	"github.com/tj/assert"
)

func TestQuery(t *testing.T) {
	url := "http://127.0.0.1:8080/v1/graphql"
	if v := os.Getenv("HASURA_URL"); v != "" {
		url = v
	}

	systemToken := os.Getenv("HASURA_SYSTEM_TOKEN")
	if systemToken == "" {
		t.Skip("HASURA_SYSTEM_TOKEN not set")
	}
	if !strings.HasPrefix(systemToken, "Bearer ") {
		systemToken = "Bearer " + systemToken
	}

	h := hasura.NewHasuraClientTo(&hasura.HasuraClientConfig{
		Url:         url,
		SystemToken: systemToken,
	})

	query := `
		query Patient {
			patient {
				id
			}
		}
	`

	result, err := h.Query(context.Background(), query, nil)
	if err != nil {
		t.Errorf("error %s", err.Error())
	}

	assert.NotNil(t, result.Get("patient"))
}
