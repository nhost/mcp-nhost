package cloud

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	ToolGetGraphqlSchemaName = "cloud-get-graphql-schema"
	//nolint:lll
	ToolGetGraphqlSchemaInstructions = `Get GraphQL schema for the Nhost Cloud allowing operations on projects and organizations. Retrieve the schema before using the tool to understand the available queries and mutations. Projects are equivalent to apps in the schema. IDs are typically uuids`
)

func (t *Tool) registerGetGraphqlSchema(mcpServer *server.MCPServer) {
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
}

func (t *Tool) handleGetGraphqlSchema(
	_ context.Context, _ mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	schema := schemaGraphql
	if t.withMutations {
		schema = schemaGraphqlWithMutations
	}

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
