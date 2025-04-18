package cloud

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/nhost/mcp-nhost/graphql"
)

//go:embed schema.graphql
var schemaGraphql string

type Tool struct {
	graphqlURL   string
	interceptors []func(ctx context.Context, req *http.Request) error
}

func NewTool(
	graphqlURL string,
	interceptors ...func(ctx context.Context, req *http.Request) error,
) *Tool {
	return &Tool{
		graphqlURL:   graphqlURL,
		interceptors: interceptors,
	}
}

func (t *Tool) Register(mcpServer *server.Server) error {
	schemaTool, err := protocol.NewTool(
		"get-graphql-schema",
		"Get GraphQL schema for the Nhost Cloud allowing operations on projects and organizations",
		GetGraphqlSchemaRequest{},
	)
	if err != nil {
		return fmt.Errorf("failed to create get-graphql-schema tool: %w", err)
	}

	mcpServer.RegisterTool(schemaTool, t.handleGetGraphqlSchema)

	queryTool, err := protocol.NewTool(
		"graphql-query",
		"Execute a GraphQL query against the Nhost Cloud to perform operations on projects and organizations",
		GraqhqlQueryRequest{}, //nolint:exhaustruct
	)
	if err != nil {
		return fmt.Errorf("failed to create get-graphql-schema tool: %w", err)
	}

	mcpServer.RegisterTool(queryTool, t.handleGraphqlQuery)

	return nil
}

type GetGraphqlSchemaRequest struct{}

func (t *Tool) handleGetGraphqlSchema(
	_ *protocol.CallToolRequest,
) (*protocol.CallToolResult, error) {
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Annotated: protocol.Annotated{
					Annotations: nil,
				},
				Type: "text",
				Text: schemaGraphql,
			},
		},
		IsError: false,
	}, nil
}

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
