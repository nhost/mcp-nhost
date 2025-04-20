package local

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/nhost/mcp-nhost/graphql"
)

const (
	ToolConfigServerQueryName = "local-config-server-query"
	//nolint:lll
	ToolConfigServerQueryInstructions = `Execute a GraphQL query against the local config server. This tool is useful to query and perform configuration changes on the local development project. Before using this tool, make sure to get the schema using the local-config-server-schema tool. To perform configuration changes this endpoint is all you need but to apply them you need to run 'nhost up' again. Ask the user for input when you need information about settings, for instance if the user asks to enable some oauth2 method and you need the client id or secret.`
)

type ConfigServerQueryRequest struct {
	Query     string `description:"graphql query to perform"      json:"query"     required:"true"`
	Variables string `description:"variables to use in the query" json:"variables" required:"false"`
}

func (t *Tool) handleConfigServerQuery(
	req *protocol.CallToolRequest,
) (*protocol.CallToolResult, error) {
	var graphReq ConfigServerQueryRequest
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
		t.configServerURL,
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
