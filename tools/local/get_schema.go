package local

import (
	"context"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/nhost/auth"
)

const (
	ToolGetGraphqlSchemaName = "local-get-graphql-schema"
	//nolint:lll
	ToolGetGraphqlSchemaInstructions = `Get GraphQL schema for an Nhost development project running locally via the Nhost CLI. This tool is useful when the user is developing a project and wants help generating code to interact with their project's Graphql schema.`
)

//nolint:lll
type GetGraphqlSchemaRequest struct {
	Role string `description:"role to use when executing queries. Default to user but make sure the user is aware" json:"role" required:"true"`
}

func (t *Tool) handleGetGraphqlSchema(
	req *protocol.CallToolRequest,
) (*protocol.CallToolResult, error) {
	var schemaReq GetGraphqlSchemaRequest
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &schemaReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	interceptors := append( //nolint:gocritic
		t.interceptors,
		auth.WithRole(schemaReq.Role),
	)

	var introspection graphql.ResponseIntrospection
	if err := graphql.Query(
		context.Background(),
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
