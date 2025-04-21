package project

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/nhost/auth"
	"github.com/nhost/mcp-nhost/tools"
)

const (
	ToolGetGraphqlSchemaName         = "project-get-graphql-schema"
	ToolGetGraphqlSchemaInstructions = `Get GraphQL schema for an Nhost project running in the Nhost Cloud.`
)

func (t *Tool) registerGetGraphqlSchemaTool(mcpServer *server.MCPServer) {
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
			mcp.Description(
				"role to use when executing queries. Default to user but make sure the user is aware",
			),
			mcp.Required(),
		),
	)
	mcpServer.AddTool(schemaTool, t.handleGetGraphqlSchema)
}

func (t *Tool) handleGetGraphqlSchema(
	ctx context.Context, req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	role, err := tools.RoleFromParams(req.Params.Arguments)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	interceptors := append( //nolint:gocritic
		t.interceptors,
		auth.WithRole(role),
	)

	var introspection graphql.ResponseIntrospection
	if err := graphql.Query(
		ctx,
		t.graphqlURL,
		graphql.IntrospectionQuery,
		nil,
		&introspection,
		interceptors...,
	); err != nil {
		return nil, fmt.Errorf("failed to query GraphQL schema: %w", err)
	}

	schema := graphql.ParseSchema(
		introspection,
		graphql.Filter{
			AllowQueries:   nil,
			AllowMutations: nil,
		},
	)

	return &mcp.CallToolResult{
		Result: mcp.Result{
			Meta: nil,
		},
		Content: []mcp.Content{
			mcp.TextContent{
				Annotated: mcp.Annotated{
					Annotations: nil,
				},
				Type: "text",
				Text: schema,
			},
		},
		IsError: false,
	}, nil
}
