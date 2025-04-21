package project

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
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
	schemaTool := mcp.NewTool(
		ToolGetGraphqlSchemaName,
		mcp.WithDescription(ToolGetGraphqlSchemaInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Get GraphQL Schema for Nhost Project running on Nhost Cloud",
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

	allowedMutations := t.allowMutations == nil || len(t.allowMutations) > 0

	queryTool := mcp.NewTool(
		ToolGraphqlQueryName,
		mcp.WithDescription(ToolGraphqlQueryInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Perform GraphQL Query on Nhost Project running on Nhost Cloud",
				ReadOnlyHint:    !allowedMutations,
				DestructiveHint: allowedMutations,
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

	return nil
}
