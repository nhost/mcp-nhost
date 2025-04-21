package local

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nhost/mcp-nhost/graphql"
)

const (
	ToolConfigServerSchemaName = "local-config-server-get-schema"
	//nolint:lll
	ToolConfigServerSchemaInstructions = `Get GraphQL schema for the local config server. This tool is useful when the user is developing a project and wants help changing the project's settings.`
)

func (t *Tool) handleConfigServerSchema(
	ctx context.Context, _ mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	var introspection graphql.ResponseIntrospection
	if err := graphql.Query(
		ctx,
		t.configServerURL,
		graphql.IntrospectionQuery,
		nil,
		&introspection,
		t.interceptors...,
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
