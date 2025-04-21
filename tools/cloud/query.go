package cloud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/tools"
)

const (
	ToolGraphqlQueryName = "cloud-graphql-query"
	//nolint:lll
	ToolGraphqlQueryInstructions = `Execute a GraphQL query against the Nhost Cloud to perform operations on projects and organizations. It also allows configuring projects hosted on Nhost Cloud. Make sure you got the schema before attempting to execute any query. If you get an error while performing a query refresh the schema in case something has changed or you did something wrong. If you get an error indicating mutations are not allowed the user may have disabled them in the server, don't retry and ask the user they need to pass --with-cloud-mutations when starting mcp-nhost to enable them. Projects are apps.`
)

func (t *Tool) registerGraphqlQuery(mcpServer *server.MCPServer) {
	t.registerGetGraphqlSchema(mcpServer)
	queryTool := mcp.NewTool(
		ToolGraphqlQueryName,
		mcp.WithDescription(ToolGraphqlQueryInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Perform GraphQL Query on Nhost Cloud Platform",
				ReadOnlyHint:    !t.withMutations,
				DestructiveHint: t.withMutations,
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
	)
	mcpServer.AddTool(queryTool, t.handleGraphqlQuery)
}

func (t *Tool) handleGraphqlQuery(
	ctx context.Context, req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	request, err := tools.QueryRequestFromParams(req.Params.Arguments)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	var resp graphql.Response[any]
	if err := graphql.Query(
		ctx,
		t.graphqlURL,
		request.Query,
		request.Variables,
		&resp,
		t.interceptors...,
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
