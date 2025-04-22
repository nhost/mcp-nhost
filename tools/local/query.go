package local

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/nhost/auth"
	"github.com/nhost/mcp-nhost/tools"
)

const (
	ToolGraphqlQueryName = "local-graphql-query"
	//nolint:lll
	ToolGraphqlQueryInstructions = `Execute a GraphQL query against an Nhost development project running locally via the Nhost CLI. This tool is useful to test queries and mutations during development. If you run into issues executing queries, retrieve the schema using the local-get-graphql-schema tool in case the schema has changed.`
)

func (t *Tool) registerGraphqlQuery(mcpServer *server.MCPServer) {
	queryTool := mcp.NewTool(
		ToolGraphqlQueryName,
		mcp.WithDescription(ToolGraphqlQueryInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Perform GraphQL Query on Nhost Development Project",
				ReadOnlyHint:    false,
				DestructiveHint: true,
				IdempotentHint:  false,
				OpenWorldHint:   true,
			},
		),
		mcp.WithString(
			"query",
			mcp.Description("graphql query to perform"),
			mcp.Required(),
		),
		mcp.WithString(
			"variables",
			mcp.Description("variables to use in the query"),
		),
		mcp.WithString(
			"role",
			mcp.Description(
				"role to use when executing queries. Default to user but make sure the user is aware. Keep in mind the schema depends on the role so if you retrieved the schema for a different role previously retrieve it for this role beforehand as it might differ", //nolint:lll
			),
			mcp.Required(),
		),
	)
	mcpServer.AddTool(queryTool, t.handleGraphqlQuery)
}

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
		nil,
		nil,
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
