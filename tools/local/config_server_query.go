package local

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/tools"
)

const (
	ToolConfigServerQueryName = "local-config-server-query"
	//nolint:lll
	ToolConfigServerQueryInstructions = `Execute a GraphQL query against the local config server. This tool is useful to query and perform configuration changes on the local development project. Before using this tool, make sure to get the schema using the local-config-server-schema tool. To perform configuration changes this endpoint is all you need but to apply them you need to run 'nhost up' again. Ask the user for input when you need information about settings, for instance if the user asks to enable some oauth2 method and you need the client id or secret.`
)

func (t *Tool) handleConfigServerQuery(
	ctx context.Context, req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	request, err := tools.QueryRequestFromParams(req.Params.Arguments)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	var resp graphql.Response[any]
	if err := graphql.Query(
		ctx,
		t.configServerURL,
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
