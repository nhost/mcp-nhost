package local

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/nhost/auth"
	"github.com/nhost/mcp-nhost/tools"
)

const (
	ToolGetGraphqlSchemaName = "local-get-graphql-schema"
	//nolint:lll
	ToolGetGraphqlSchemaInstructions = `Get GraphQL schema for an Nhost development project running locally via the Nhost CLI. This tool is useful when the user is developing a project and wants help generating code to interact with their project's Graphql schema.`
)

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
