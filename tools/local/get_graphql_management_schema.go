package local

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nhost/mcp-nhost/nhost/graphql"
)

const (
	ToolGetGraphqlManagementSchemaName         = "local-get-management-graphql-schema"
	ToolGetGraphqlManagementSchemaInstructions = `
		Get GraphQL's management schema for an Nhost development project running locally via the Nhost
		CLI. This tool is useful to properly understand how manage hasura metadata, migrations,
		permissions, remote schemas, etc.

		Use it before attempting to use ` + ToolManageGraphqlName
)

func (t *Tool) registerGetGraphqlManagementSchema(mcpServer *server.MCPServer) {
	schemaTool := mcp.NewTool(
		ToolGetGraphqlManagementSchemaName,
		mcp.WithDescription(ToolGetGraphqlManagementSchemaInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Get GraphQL's Management Schema for Nhost Development Project",
				ReadOnlyHint:    true,
				DestructiveHint: false,
				IdempotentHint:  true,
				OpenWorldHint:   true,
			},
		),
	)
	mcpServer.AddTool(schemaTool, t.handleGetGraphqlManagementSchema)
}

func (t *Tool) handleGetGraphqlManagementSchema(
	_ context.Context, _ mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
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
				Text: graphql.Schema,
			},
		},
		IsError: false,
	}, nil
}
