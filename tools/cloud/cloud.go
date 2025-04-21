package cloud

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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

func (t *Tool) Register(mcpServer *server.MCPServer) error {
	schemaTool := mcp.NewTool(
		ToolGetGraphqlSchemaName,
		mcp.WithDescription(ToolGetGraphqlSchemaInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Get GraphQL Schema for Nhost Cloud Platform",
				ReadOnlyHint:    true,
				DestructiveHint: false,
				IdempotentHint:  true,
				OpenWorldHint:   true,
			},
		),
	)
	mcpServer.AddTool(schemaTool, t.handleGetGraphqlSchema)

	queryTool := mcp.NewTool(
		ToolGraphqlQueryName,
		mcp.WithDescription(ToolGraphqlQueryInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Perform GraphQL Query on Nhost Cloud Platform",
				ReadOnlyHint:    !t.withMutations,
				DestructiveHint: t.withMutations,
				IdempotentHint:  false,
				OpenWorldHint:   true,
			},
		),
		mcp.WithString(
			"query",
			mcp.Description("graphql query to perform"),
			mcp.Required(),
		),
		mcp.WithString(
			"variables",
			mcp.Description("variables to use in the query"),
		),
	)
	mcpServer.AddTool(queryTool, t.handleGraphqlQuery)

	return nil
}
