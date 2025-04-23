package local

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	ToolManageGraphqlName         = "local-manage-graphql"
	ToolManageGraphqlInstructions = `
		Query GraphQL's management endpoints on an Nhost development project running locally via
		the Nhost CLI. This tool is useful to manage hasura metadata, migrations, permissions,
		remote schemas, etc. Use this tool to manage the underlying database schema, specially if users
		ask to modify the project's schema as that can't be done via the GraphQL API.
		1. When ask to modify the schema:
		   1. Get the schema from the database first for reference
		   2. Always create database migrations using /apis/migrate endpoint
		   3. Don't forget to provide a down migration
		   4. Provide a pg_track_table in the same request if a new table is created
		2. When creating foreign keys, either on new tables or existing ones, make sure you track the
		   relationships
		3. When managing permissions admin always has full access so you don't need to manage
		   this particularrole
		`
)

func (t *Tool) registerManageGraphql(mcpServer *server.MCPServer) {
	schemaTool := mcp.NewTool(
		ToolManageGraphqlName,
		mcp.WithDescription(ToolManageGraphqlInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Manage GraphQL's Metadata on an Nhost Development Project",
				ReadOnlyHint:    false,
				DestructiveHint: true,
				IdempotentHint:  true,
				OpenWorldHint:   true,
			},
		),
		mcp.WithString(
			"endpoint",
			mcp.Description("The GraphQL management endpoint to query. Use https://local.hasura.local.nhost.run as base URL"),
			mcp.Required(),
		),
		mcp.WithString(
			"body",
			mcp.Description("The body for the HTTP request"),
			mcp.Required(),
		),
		mcp.WithString(
			"method",
			mcp.Description("The HTTP method to use"),
			mcp.DefaultString("POST"),
			mcp.Enum("GET", "POST", "PUT", "DELETE"),
		),
	)
	mcpServer.AddTool(schemaTool, t.handleManageGraphql)
}

type httpResponse struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

func genericQuery(
	ctx context.Context,
	endpoint string,
	body string,
	method string,
	headers http.Header,
	interceptors []func(ctx context.Context, req *http.Request) error,
) (httpResponse, error) {
	request, err := http.NewRequestWithContext(ctx, method, endpoint, strings.NewReader(body))
	if err != nil {
		return httpResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	for key, values := range headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	for _, interceptor := range interceptors {
		if err := interceptor(ctx, request); err != nil {
			return httpResponse{}, fmt.Errorf("failed to execute interceptor: %w", err)
		}
	}

	client := &http.Client{} //nolint: exhaustruct
	response, err := client.Do(request)
	if err != nil {
		return httpResponse{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)

	return httpResponse{
		StatusCode: response.StatusCode,
		Body:       string(b),
	}, nil
}

func (t *Tool) handleManageGraphqlArguments(
	arguments map[string]any,
) (string, string, string, http.Header, error) {
	endpoint, ok := arguments["endpoint"].(string)
	if !ok {
		return "", "", "", nil, fmt.Errorf("missing endpoint")
	}

	body, ok := arguments["body"].(string)
	if !ok {
		return "", "", "", nil, fmt.Errorf("missing body")
	}

	method, ok := arguments["method"].(string)
	if !ok {
		return "", "", "", nil, fmt.Errorf("missing method")
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("Accept", "application/json")

	return endpoint, body, method, headers, nil
}

func (t *Tool) handleManageGraphql(
	ctx context.Context, req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	endpoint, body, method, headers, err := t.handleManageGraphqlArguments(req.Params.Arguments)

	response, err := genericQuery(ctx, endpoint, body, method, headers, t.interceptors)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	b, err := json.Marshal(response)
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
