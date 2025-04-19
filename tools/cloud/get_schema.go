package cloud

import "github.com/ThinkInAIXYZ/go-mcp/protocol"

const (
	toolGetGraphqlSchemaName = "cloud-get-graphql-schema"
	//nolint:lll
	toolGetGraphqlSchemaInstructions = `Get GraphQL schema for the Nhost Cloud allowing operations on projects and organizations.`
)

type GetGraphqlSchemaRequest struct{}

func (t *Tool) handleGetGraphqlSchema(
	_ *protocol.CallToolRequest,
) (*protocol.CallToolResult, error) {
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Annotated: protocol.Annotated{
					Annotations: nil,
				},
				Type: "text",
				Text: schemaGraphql,
			},
		},
		IsError: false,
	}, nil
}
