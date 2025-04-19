package cloud

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

//go:embed schema.graphql
var schemaGraphql string

type Tool struct {
	graphqlURL   string
	interceptors []func(ctx context.Context, req *http.Request) error
}

func NewTool(
	graphqlURL string,
	interceptors ...func(ctx context.Context, req *http.Request) error,
) *Tool {
	return &Tool{
		graphqlURL:   graphqlURL,
		interceptors: interceptors,
	}
}

func (t *Tool) Register(mcpServer *server.Server) error {
	schemaTool, err := protocol.NewTool(
		toolGetGraphqlSchemaName,
		toolGetGraphqlSchemaInstructions,
		GetGraphqlSchemaRequest{},
	)
	if err != nil {
		return fmt.Errorf("failed to create %s tool: %w", toolGetGraphqlSchemaName, err)
	}

	mcpServer.RegisterTool(schemaTool, t.handleGetGraphqlSchema)

	queryTool, err := protocol.NewTool(
		toolGraphqlQueryName,
		toolGraphqlQueryInstructions,
		GraqhqlQueryRequest{}, //nolint:exhaustruct
	)
	if err != nil {
		return fmt.Errorf("failed to create %s tool: %w", toolGraphqlQueryName, err)
	}

	mcpServer.RegisterTool(queryTool, t.handleGraphqlQuery)

	return nil
}
