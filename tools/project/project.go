package project

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

type Tool struct {
	graphqlURL     string
	allowQueries   []string
	allowMutations []string
	interceptors   []func(ctx context.Context, req *http.Request) error
}

func NewTool(
	graphqlURL string,
	allowQueries []string,
	allowMutations []string,
	interceptors ...func(ctx context.Context, req *http.Request) error,
) *Tool {
	return &Tool{
		graphqlURL:     graphqlURL,
		allowQueries:   allowQueries,
		allowMutations: allowMutations,
		interceptors:   interceptors,
	}
}

func (t *Tool) Register(mcpServer *server.MCPServer) error {
	t.registerGetGraphqlSchemaTool(mcpServer)
	t.registerGraphqlQuery(mcpServer)

	return nil
}
