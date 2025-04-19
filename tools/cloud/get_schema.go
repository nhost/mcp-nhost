package cloud

import "github.com/ThinkInAIXYZ/go-mcp/protocol"

const (
	ToolGetGraphqlSchemaName = "cloud-get-graphql-schema"
	//nolint:lll
	ToolGetGraphqlSchemaInstructions = `Get GraphQL schema for the Nhost Cloud allowing operations on projects and organizations. Retrieve the schema before using the tool to understand the available queries and mutations. Projects are equivalent to apps in the schema. IDs are typically uuids`
)

type GetGraphqlSchemaRequest struct{}

func (t *Tool) handleGetGraphqlSchema(
	_ *protocol.CallToolRequest,
) (*protocol.CallToolResult, error) {
	schema := schemaGraphql
	if t.withMutations {
		schema = schemaGraphqlWithMutations
	}

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
