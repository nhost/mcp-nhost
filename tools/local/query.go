package local

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/nhost/auth"
)

const (
	ToolGraphqlQueryName = "local-graphql-query"
	//nolint:lll
	ToolGraphqlQueryInstructions = `Execute a GraphQL query against an Nhost development project running locally via the Nhost CLI. This tool is useful to test queries and mutations duding development. If you run into issues executing queries, retrieve the schema using the local-get-graphql-schema tool in case the schema has changed.`
)

//nolint:lll
type GraphqlQueryRequest struct {
	Query     string `description:"graphql query to perform"      json:"query"     required:"true"`
	Variables string `description:"variables to use in the query" json:"variables" required:"false"`

	Role string `description:"role to use when executing queries. Default to user but make sure the user is aware. Keep in mind the schema depends on the role so if you retrieved the schema for a different role previously retrieve it for this role beforhand as it might differ" json:"role" required:"true"`
}

func (t *Tool) handleGraphqlQuery(req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var graphReq GraphqlQueryRequest
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &graphReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	var variables map[string]any
	if graphReq.Variables != "" {
		if err := json.Unmarshal([]byte(graphReq.Variables), &variables); err != nil {
			return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
		}
	}

	interceptors := append( //nolint:gocritic
		t.interceptors,
		auth.WithRole(graphReq.Role),
	)

	var resp graphql.Response[any]
	if err := graphql.Query(
		context.Background(),
		t.graphqlURL,
		graphReq.Query,
		variables,
		&resp,
		interceptors...,
	); err != nil {
		return nil, fmt.Errorf("failed to query graphql endpoint: %w", err)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Annotated: protocol.Annotated{
					Annotations: nil,
				},
				Type: "text",
				Text: string(b),
			},
		},
		IsError: false,
	}, nil
}
