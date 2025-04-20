package local

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

type Tool struct {
	graphqlURL      string
	configServerURL string
	interceptors    []func(ctx context.Context, req *http.Request) error
}

func NewTool(
	graphqlURL string,
	configServerURL string,
	interceptors ...func(ctx context.Context, req *http.Request) error,
) *Tool {
	return &Tool{
		graphqlURL:      graphqlURL,
		configServerURL: configServerURL,
		interceptors:    interceptors,
	}
}

func (t *Tool) Register(mcpServer *server.Server) error {
	schemaTool, err := protocol.NewTool(
		ToolGetGraphqlSchemaName,
		ToolGetGraphqlSchemaInstructions,
		GetGraphqlSchemaRequest{}, //nolint:exhaustruct
	)
	if err != nil {
		return fmt.Errorf("failed to create %s tool: %w", ToolGetGraphqlSchemaName, err)
	}

	mcpServer.RegisterTool(schemaTool, t.handleGetGraphqlSchema)

	queryTool, err := protocol.NewTool(
		ToolGraphqlQueryName,
		ToolGraphqlQueryInstructions,
		GraphqlQueryRequest{}, //nolint:exhaustruct
	)
	if err != nil {
		return fmt.Errorf("failed to create %s tool: %w", ToolGraphqlQueryName, err)
	}

	mcpServer.RegisterTool(queryTool, t.handleGraphqlQuery)

	configServerSchemaTool, err := protocol.NewTool(
		ToolConfigServerSchemaName,
		ToolConfigServerSchemaInstructions,
		ConfigServerSchemaRequest{},
	)
	if err != nil {
		return fmt.Errorf("failed to create %s tool: %w", ToolConfigServerSchemaName, err)
	}

	mcpServer.RegisterTool(configServerSchemaTool, t.handleConfigServerSchema)

	configServerQueryTool, err := protocol.NewTool(
		ToolConfigServerQueryName,
		ToolConfigServerQueryInstructions,
		ConfigServerQueryRequest{}, //nolint:exhaustruct
	)
	if err != nil {
		return fmt.Errorf("failed to create %s tool: %w", ToolConfigServerQueryName, err)
	}

	mcpServer.RegisterTool(configServerQueryTool, t.handleConfigServerQuery)

	return nil
}
