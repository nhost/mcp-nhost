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

//go:embed schema-with-mutations.graphql
var schemaGraphqlWithMutations string

type Tool struct {
	graphqlURL    string
	withMutations bool
	interceptors  []func(ctx context.Context, req *http.Request) error
}

func NewTool(
	graphqlURL string,
	withMutations bool,
	interceptors ...func(ctx context.Context, req *http.Request) error,
) *Tool {
	return &Tool{
		graphqlURL:    graphqlURL,
		withMutations: withMutations,
		interceptors:  interceptors,
	}
}

func (t *Tool) Register(mcpServer *server.Server) error {
	schemaTool, err := protocol.NewTool(
		ToolGetGraphqlSchemaName,
		ToolGetGraphqlSchemaInstructions,
		GetGraphqlSchemaRequest{},
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

	return nil
}
