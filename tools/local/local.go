package local

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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

func (t *Tool) Register(mcpServer *server.MCPServer) error {
	schemaTool := mcp.NewTool(
		ToolGetGraphqlSchemaName,
		mcp.WithDescription(ToolGetGraphqlSchemaInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Get GraphQL Schema for Nhost Development Project",
				ReadOnlyHint:    true,
				DestructiveHint: false,
				IdempotentHint:  true,
				OpenWorldHint:   true,
			},
		),
		mcp.WithString(
			"role",
			mcp.Description("role to use when executing queries. Default to user but make sure the user is aware"),
			mcp.Required(),
		),
	)
	mcpServer.AddTool(schemaTool, t.handleGetGraphqlSchema)

	queryTool := mcp.NewTool(
		ToolGraphqlQueryName,
		mcp.WithDescription(ToolGraphqlQueryInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Perform GraphQL Query on Nhost Development Project",
				ReadOnlyHint:    false,
				DestructiveHint: true,
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
		mcp.WithString(
			"role",
			mcp.Description("role to use when executing queries. Default to user but make sure the user is aware. Keep in mind the schema depends on the role so if you retrieved the schema for a different role previously retrieve it for this role beforehand as it might differ"),
			mcp.Required(),
		),
	)
	mcpServer.AddTool(queryTool, t.handleGraphqlQuery)

	configServerSchemaTool := mcp.NewTool(
		ToolConfigServerSchemaName,
		mcp.WithDescription(ToolConfigServerSchemaInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Get GraphQL Schema for Nhost Config Server",
				ReadOnlyHint:    true,
				DestructiveHint: false,
				IdempotentHint:  true,
				OpenWorldHint:   true,
			},
		),
	)
	mcpServer.AddTool(configServerSchemaTool, t.handleConfigServerSchema)

	configServerQueryTool := mcp.NewTool(
		ToolConfigServerQueryName,
		mcp.WithDescription(ToolConfigServerQueryInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Perform GraphQL Query on Nhost Config Server",
				ReadOnlyHint:    false,
				DestructiveHint: true,
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
	mcpServer.AddTool(configServerQueryTool, t.handleConfigServerQuery)

	return nil
}
