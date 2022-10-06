package hasura

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

// Connect with an Hasura GraphQL API
type HasuraClient interface {
	Query(query string, vars *Variables) (*QueryResult, error)
	Mutate(mutation string, vars *Variables) (*QueryResult, error)
}

type HasuraClientConfig struct {
	Url         string
	SystemToken string
}

type hasura struct {
	client      graphql.Client
	systemToken string
}

// Creates an HasuraClient
func NewHasuraClient(GraphQLServerURL, SystemToken string) HasuraClient {
	return NewHasuraClientTo(&HasuraClientConfig{
		Url:         GraphQLServerURL,
		SystemToken: fmt.Sprintf("Bearer %s", SystemToken),
	})
}

// Creates an HasuraClient with provided config.
func NewHasuraClientTo(config *HasuraClientConfig) HasuraClient {
	return &hasura{
		client:      *graphql.NewClient(config.Url),
		systemToken: config.SystemToken,
	}
}

// Perform queries, to retrieve data from database.
func (h *hasura) Query(query string, vars *Variables) (*QueryResult, error) {
	return h.run(query, vars)
}

// Perform mutations to change database state.
func (h *hasura) Mutate(mutation string, vars *Variables) (*QueryResult, error) {
	return h.run(mutation, vars)
}

func (h *hasura) run(cmd string, vars *Variables) (*QueryResult, error) {
	req := graphql.NewRequest(cmd)
	req.Header.Add("x-hasura-role", "system")
	req.Header.Add("Authorization", h.systemToken)
	req.Header.Set("Cache-Control", "no-cache")

	if vars != nil {
		for key, value := range *vars {
			req.Var(key, value)
		}
	}

	var response map[string]interface{}

	err := h.client.Run(context.Background(), req, &response)
	if err != nil {
		return nil, fmt.Errorf("[HasuraClient][run] error: %w", err)
	}

	result := &QueryResult{
		Data: response,
	}

	return result, nil
}
