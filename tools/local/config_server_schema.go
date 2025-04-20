package local

import (
	"context"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/nhost/mcp-nhost/graphql"
)

const (
	ToolConfigServerSchemaName = "local-config-server-schema"
	//nolint:lll
	ToolConfigServerSchemaInstructions = `Get GraphQL schema for the local config server. This tool is useful when the user is developing a project and wants help changing the project's settings.`
)

type ConfigServerSchemaRequest struct{}

func (t *Tool) handleConfigServerSchema(_ *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var introspection graphql.ResponseIntrospection
	if err := graphql.Query(
		context.Background(),
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

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Annotated: protocol.Annotated{
					Annotations: nil,
				},
				Type: "text",
				Text: schema,
			},
		},
		IsError: false,
	}, nil
}
