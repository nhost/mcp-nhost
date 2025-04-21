package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/nhost/auth"
	"github.com/nhost/mcp-nhost/tools"
)

const (
	ToolGraphqlQueryName = "project-graphql-query"
	//nolint:lll
	ToolGraphqlQueryInstructions = `Execute a GraphQL query against a Nhost project running in the Nhost Cloud. This tool is useful to query and mutate live data running on an online projec. If you run into issues executing queries, retrieve the schema using the project-get-graphql-schema tool in case the schema has changed. If you get an error indicating the query or mutation is not allowed the user may have disabled them in the server, don't retry and tell the user they need to enable them when starting mcp-nhost`
)

func (t *Tool) handleGraphqlQuery(
	ctx context.Context, req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	request, err := tools.QueryRequestWithRoleFromParams(req.Params.Arguments)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	interceptors := append( //nolint:gocritic
		t.interceptors,
		auth.WithRole(request.Role),
	)

	var resp graphql.Response[any]
	if err := graphql.Query(
		ctx,
		t.graphqlURL,
		request.Query,
		request.Variables,
		&resp,
		interceptors...,
	); err != nil {
		return nil, fmt.Errorf("failed to query graphql endpoint: %w", err)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
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
				Text: string(b),
			},
		},
		IsError: false,
	}, nil
}
