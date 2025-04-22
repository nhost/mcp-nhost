package project

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/nhost/auth"
	"github.com/nhost/mcp-nhost/tools"
)

const (
	ToolGetGraphqlSchemaName         = "project-get-graphql-schema"
	ToolGetGraphqlSchemaInstructions = `Get GraphQL schema for an Nhost project running in the Nhost Cloud.`
)

var ErrNotFound = errors.New("not found")

func (t *Tool) registerGetGraphqlSchemaTool(mcpServer *server.MCPServer, projects string) {
	schemaTool := mcp.NewTool(
		ToolGetGraphqlSchemaName,
		mcp.WithDescription(ToolGetGraphqlSchemaInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Get GraphQL Schema for Nhost Project running on Nhost Cloud",
				ReadOnlyHint:    true,
				DestructiveHint: false,
				IdempotentHint:  true,
				OpenWorldHint:   true,
			},
		),
		mcp.WithString(
			"role",
			mcp.Description(
				"role to use when executing queries. Default to user but make sure the user is aware",
			),
			mcp.Required(),
		),
		mcp.WithString(
			"projectSubdomain",
			mcp.Description(
				fmt.Sprintf(
					"Project to get the GraphQL schema for. Must be one of %s, otherwise you don't have access to it. You can use cloud-* tools to resolve subdomains and map them to names", //nolint:lll
					projects,
				),
			),
			mcp.Required(),
		),
	)
	mcpServer.AddTool(schemaTool, t.handleGetGraphqlSchema)
}

func (t *Tool) handleGetGraphqlSchema(
	ctx context.Context, req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	role, err := tools.RoleFromParams(req.Params.Arguments)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	projectSubdomain, err := tools.ProjectFromParams(req.Params.Arguments)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	project, ok := t.projects[projectSubdomain]
	if !ok {
		return nil,
			errors.New("this project is not configured to be accessed by an LLM") //nolint:goerr113
	}

	interceptors := []func(ctx context.Context, req *http.Request) error{
		project.authInterceptor,
		auth.WithRole(role),
	}

	var introspection graphql.ResponseIntrospection
	if err := graphql.Query(
		ctx,
		project.graphqlURL,
		graphql.IntrospectionQuery,
		nil,
		&introspection,
		nil,
		nil,
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
				Text: schema,
			},
		},
		IsError: false,
	}, nil
}
