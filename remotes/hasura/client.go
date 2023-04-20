package hasura

import (
	"context"
	"fmt"

	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
	"github.com/machinebox/graphql"
)

// Connect with an Hasura GraphQL API
type HasuraClient interface {
	Query(ctx context.Context, query string, vars *Variables) (*QueryResult, error)
	Mutate(ctx context.Context, mutation string, vars *Variables) (*QueryResult, error)
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

func (h *hasura) tags(method, query string) map[string]interface{} {
	return map[string]interface{}{
		"method": method,
		"query":  query,
	}
}

// Perform queries, to retrieve data from database.
func (h *hasura) Query(ctx context.Context, query string, vars *Variables) (*QueryResult, error) {
	_, span := jaeger.SpanTrace(ctx, "hasura.query", h.tags("query", query))
	defer span.Finish()
	return h.run(query, vars)
}

// Perform mutations to change database state.
func (h *hasura) Mutate(ctx context.Context, mutation string, vars *Variables) (*QueryResult, error) {
	_, span := jaeger.SpanTrace(ctx, "hasura.mutation", h.tags("query", mutation))
	defer span.Finish()
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
