package cloud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/nhost/mcp-nhost/graphql"
)

const (
	toolGraphqlQueryName = "cloud-graphql-query"
	//nolint:lll
	toolGraphqlQueryInstructions = `Execute a GraphQL query against the Nhost Cloud to perform operations on projects and organizations. It also allows configuring projects hosted on Nhost Cloud.`
)

type GraqhqlQueryRequest struct {
	Query     string `description:"graphql query to perform"      json:"query"     required:"true"`
	Variables string `description:"variables to use in the query" json:"variables" required:"false"`
}

func (t *Tool) handleGraphqlQuery(req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var graphReq GraqhqlQueryRequest
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &graphReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	var variables map[string]any
	if graphReq.Variables != "" {
		if err := json.Unmarshal([]byte(graphReq.Variables), &variables); err != nil {
			return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
		}
	}

	var resp graphql.Response[any]
	if err := graphql.Query(
		context.Background(),
		t.graphqlURL,
		graphReq.Query,
		variables,
		&resp,
		t.interceptors...,
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
