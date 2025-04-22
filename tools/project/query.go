package project

import (
	"context"
	"encoding/json"
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
	ToolGraphqlQueryName = "project-graphql-query"
	//nolint:lll
	ToolGraphqlQueryInstructions = `Execute a GraphQL query against a Nhost project running in the Nhost Cloud. This tool is useful to query and mutate live data running on an online projec. If you run into issues executing queries, retrieve the schema using the project-get-graphql-schema tool in case the schema has changed. If you get an error indicating the query or mutation is not allowed the user may have disabled them in the server, don't retry and tell the user they need to enable them when starting mcp-nhost`
)

func (t *Tool) registerGraphqlQuery(mcpServer *server.MCPServer, projects string) {
	allowedMutations := false

	for _, proj := range t.projects {
		if proj.allowMutations == nil || len(proj.allowMutations) > 0 {
			allowedMutations = true
			break
		}
	}

	queryTool := mcp.NewTool(
		ToolGraphqlQueryName,
		mcp.WithDescription(ToolGraphqlQueryInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Perform GraphQL Query on Nhost Project running on Nhost Cloud",
				ReadOnlyHint:    !allowedMutations,
				DestructiveHint: allowedMutations,
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
			"projectSubdomain",
			mcp.Description(
				fmt.Sprintf(
					"Project to get the GraphQL schema for. Must be one of %s, otherwise you don't have access to it. You can use cloud-* tools to resolve subdomains and map them to names", //nolint:lll
					projects,
				),
			),
			mcp.Required(),
		),
		mcp.WithString(
			"role",
			mcp.Description(
				"role to use when executing queries. Default to user but make sure the user is aware. Keep in mind the schema depends on the role so if you retrieved the schema for a different role previously retrieve it for this role beforehand as it might differ", //nolint:lll
			),
			mcp.Required(),
		),
		mcp.WithString(
			"userId",
			mcp.Description(
				"Overrides X-Hasura-User-Id in the GraphQL query/mutation. Credentials must allow it (i.e. admin secret must be in use)", //nolint:lll
			),
		),
	)
	mcpServer.AddTool(queryTool, t.handleGraphqlQuery)
}

func (t *Tool) handleGraphqlQueryArgs(
	args map[string]any,
) (tools.GraphqlQueryWithRoleRequest, string, string, error) {
	request, err := tools.QueryRequestWithRoleFromParams(args)
	if err != nil {
		return request, "", "", err //nolint:wrapcheck
	}

	projectSubdomain, err := tools.ProjectFromParams(args)
	if err != nil {
		return request, "", "", err //nolint:wrapcheck
	}

	userID, err := tools.FromParams[string](args, "userId")
	if err != nil {
		return request, "", "", err
	}

	return request, projectSubdomain, userID, nil
}

func (t *Tool) handleGraphqlQuery(
	ctx context.Context, req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	request, projectSubdomain, userID, err := t.handleGraphqlQueryArgs(req.Params.Arguments)
	if err != nil {
		return nil, err
	}

	project, ok := t.projects[projectSubdomain]
	if !ok {
		return nil,
			errors.New("this project is not configured to be accessed by an LLM") //nolint:goerr113
	}

	interceptors := []func(ctx context.Context, req *http.Request) error{
		project.authInterceptor,
		auth.WithRole(request.Role),
	}

	if userID != "" {
		interceptors = append(interceptors, auth.WithUserID(userID))
	}

	var resp graphql.Response[any]
	if err := graphql.Query(
		ctx,
		project.graphqlURL,
		request.Query,
		request.Variables,
		&resp,
		project.allowQueries,
		project.allowMutations,
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
