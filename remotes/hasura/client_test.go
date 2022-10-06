package hasura_test

import (
	"testing"

	"github.com/joaopandolfi/blackwhale/remotes/hasura"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	h := hasura.NewHasuraClientTo(&hasura.HasuraClientConfig{
		Url:         "http://localhost:8080/v1/graphql",
		SystemToken: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJicm9rZXIiOnsieC1oYXN1cmEtYWxsb3dlZC1yb2xlcyI6WyJzeXN0ZW0iLCJ1c2VyIl0sIngtaGFzdXJhLWRlZmF1bHQtcm9sZSI6InVzZXIiLCJ4LWhhc3VyYS11c2VyLWlkIjoiZjc2YzIwY2UtYjZlZS00MjY2LTk3ODgtNWM5MTUxZWFiYjZhIn0sImV4cCI6MTY2Mzc5MTQ1MSwiaWQiOiJmNzZjMjBjZS1iNmVlLTQyNjYtOTc4OC01YzkxNTFlYWJiNmEiLCJpbnN0aXR1dGlvbiI6InNhdXJvbiIsInBlcm1pc3Npb24iOiJzeXN0ZW07dXNlciJ9.aAWGeyvyHEm526_3pxBf3zq6niCUmPqL9IzmQZKuY4c",
	})

	query := `
		query Patient {
			patient {
				id
			}
		}
	`

	result, err := h.Query(query, nil)
	if err != nil {
		t.Errorf("error %s", err.Error())
	}

	assert.NotNil(t, result.Get("patient"))
}
