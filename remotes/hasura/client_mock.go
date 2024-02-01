package hasura

import (
	"context"
	"fmt"

	"github.com/joaopandolfi/blackwhale/utils"
)

type HasuraClientMockResponse struct {
	Query  string
	Vars   *Variables
	Result *QueryResult
	Error  error
}

type HasuraClientMock struct {
	Responses []HasuraClientMockResponse
	Handler   func(query string, vars *Variables) (*QueryResult, error)
}

func NewHasuraClientMock(
	responses []HasuraClientMockResponse,
	handler func(query string, vars *Variables) (*QueryResult, error),
) HasuraClient {
	return &HasuraClientMock{
		Responses: responses,
		Handler:   handler,
	}
}

func (c *HasuraClientMock) Query(ctx context.Context, query string, vars *Variables) (*QueryResult, error) {
	return c.handle(query, vars)
}

func (c *HasuraClientMock) Mutate(ctx context.Context, mutation string, vars *Variables) (*QueryResult, error) {
	return c.handle(mutation, vars)
}

func (c *HasuraClientMock) handle(query string, vars *Variables) (*QueryResult, error) {
	if c.Handler == nil {
		return defaultHandler(query, vars, c.Responses)
	}

	return c.Handler(query, vars)
}

func defaultHandler(query string, vars *Variables, responses []HasuraClientMockResponse) (*QueryResult, error) {
	for _, r := range responses {
		if r.Query == query && utils.Equals(*r.Vars, *vars) {
			return r.Result, r.Error
		}
	}

	return nil, fmt.Errorf("mock: not found")
}
