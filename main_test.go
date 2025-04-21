package main

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nhost/mcp-nhost/cmd/start"
	"github.com/nhost/mcp-nhost/tools/cloud"
	"github.com/nhost/mcp-nhost/tools/local"
	"github.com/nhost/mcp-nhost/tools/project"
)

func TestStart(t *testing.T) { //nolint:cyclop,maintidx
	t.Parallel()

	cmd := Command()

	buf := bytes.NewBuffer(nil)
	cmd.Writer = buf

	go func() {
		if err := cmd.Run(
			context.Background(),
			[]string{
				"main",
				"start",
				"--bind=:9000",
				"--nhost-pat=asdasd",
				"--with-cloud-mutations",
				"--project-subdomain=fake-subdomain",
				"--project-region=fake-region",
				"--project-admin-secret=fake-admin-secret",
				"--project-allow-queries=*",
				"--project-allow-mutations=*",
			},
		); err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)

	transportClient, err := transport.NewSSE("http://localhost:9000/sse")
	if err != nil {
		t.Fatalf("failed to create transport client: %v", err)
	}

	mcpClient := client.NewClient(transportClient)

	if err := mcpClient.Start(context.Background()); err != nil {
		t.Fatalf("failed to start mcp client: %v", err)
	}
	defer mcpClient.Close()

	initRequest := mcp.InitializeRequest{} //nolint:exhaustruct
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}
	res, err := mcpClient.Initialize(
		context.Background(),
		initRequest,
	)
	if err != nil {
		t.Fatalf("failed to initialize mcp client: %v", err)
	}

	if diff := cmp.Diff(
		res,
		//nolint:tagalign
		&mcp.InitializeResult{
			ProtocolVersion: "2024-11-05",
			Capabilities: mcp.ServerCapabilities{
				Experimental: nil,
				Logging:      nil,
				Prompts: &struct {
					ListChanged bool "json:\"listChanged,omitempty\""
				}{
					ListChanged: true,
				},
				Resources: &struct {
					Subscribe   bool "json:\"subscribe,omitempty\""
					ListChanged bool "json:\"listChanged,omitempty\""
				}{
					Subscribe:   true,
					ListChanged: true,
				},
				Tools: &struct {
					ListChanged bool "json:\"listChanged,omitempty\""
				}{
					ListChanged: true,
				},
			},
			ServerInfo: mcp.Implementation{
				Name:    "nhost-mcp",
				Version: "dev",
			},
			Instructions: start.ServerInstructions,
			Result: mcp.Result{
				Meta: nil,
			},
		},
	); diff != "" {
		t.Errorf("ServerInfo mismatch (-want +got):\n%s", diff)
	}

	tools, err := mcpClient.ListTools(
		context.Background(),
		mcp.ListToolsRequest{}, //nolint:exhaustruct
	)
	if err != nil {
		t.Fatalf("failed to list tools: %v", err)
	}

	if diff := cmp.Diff(
		tools,
		//nolint:exhaustruct,lll
		&mcp.ListToolsResult{
			Tools: []mcp.Tool{
				{
					Name:        "cloud-get-graphql-schema",
					Description: cloud.ToolGetGraphqlSchemaInstructions,
					InputSchema: mcp.ToolInputSchema{
						Type:       "object",
						Properties: nil,
						Required:   nil,
					},
				},
				{
					Name:        "cloud-graphql-query",
					Description: cloud.ToolGraphqlQueryInstructions,
					InputSchema: mcp.ToolInputSchema{
						Type: "object",
						Properties: map[string]any{
							"query": map[string]any{
								"description": "graphql query to perform",
								"type":        "string",
							},
							"variables": map[string]any{
								"description": "variables to use in the query",
								"type":        "string",
							},
						},
						Required: []string{"query"},
					},
				},
				{
					Name:        "local-config-server-query",
					Description: local.ToolConfigServerQueryInstructions,
					InputSchema: mcp.ToolInputSchema{
						Type: "object",
						Properties: map[string]any{
							"query": map[string]any{
								"description": "graphql query to perform",
								"type":        "string",
							},
							"variables": map[string]any{
								"description": "variables to use in the query",
								"type":        "string",
							},
						},
						Required: []string{"query"},
					},
				},
				{
					Name:        "local-config-server-schema",
					Description: local.ToolConfigServerSchemaInstructions,
					InputSchema: mcp.ToolInputSchema{
						Type:       "object",
						Properties: nil,
						Required:   nil,
					},
				},
				{
					Name:        "local-get-graphql-schema",
					Description: local.ToolGetGraphqlSchemaInstructions,
					InputSchema: mcp.ToolInputSchema{
						Type: "object",
						Properties: map[string]any{
							"role": map[string]any{
								"description": "role to use when executing queries. Default to user but make sure the user is aware",
								"type":        "string",
							},
						},
						Required: []string{"role"},
					},
				},
				{
					Name:        "local-graphql-query",
					Description: local.ToolGraphqlQueryInstructions,
					InputSchema: mcp.ToolInputSchema{
						Type: "object",
						Properties: map[string]any{
							"query": map[string]any{
								"description": "graphql query to perform",
								"type":        "string",
							},
							"role": map[string]any{
								"description": "role to use when executing queries. Default to user but make sure the user is aware. Keep in mind the schema depends on the role so if you retrieved the schema for a different role previously retrieve it for this role beforehand as it might differ",
								"type":        "string",
							},
							"variables": map[string]any{
								"description": "variables to use in the query",
								"type":        "string",
							},
						},
						Required: []string{"query", "role"},
					},
				},
				{
					Name:        "project-get-graphql-schema",
					Description: project.ToolGetGraphqlSchemaInstructions,
					InputSchema: mcp.ToolInputSchema{
						Type: "object",
						Properties: map[string]any{
							"role": map[string]any{
								"description": "role to use when executing queries. Default to user but make sure the user is aware",
								"type":        "string",
							},
						},
						Required: []string{"role"},
					},
				},
				{
					Name:        "project-graphql-query",
					Description: project.ToolGraphqlQueryInstructions,
					InputSchema: mcp.ToolInputSchema{
						Type: "object",
						Properties: map[string]any{
							"query": map[string]any{
								"description": "graphql query to perform",
								"type":        "string",
							},
							"role": map[string]any{
								"description": "role to use when executing queries. Default to user but make sure the user is aware. Keep in mind the schema depends on the role so if you retrieved the schema for a different role previously retrieve it for this role beforehand as it might differ",
								"type":        "string",
							},
							"variables": map[string]any{
								"description": "variables to use in the query",
								"type":        "string",
							},
						},
						Required: []string{"query", "role"},
					},
				},
			},
		},
		cmpopts.SortSlices(func(a, b mcp.Tool) bool {
			return a.Name < b.Name
		}),
	); diff != "" {
		t.Errorf("ListToolsResult mismatch (-want +got):\n%s", diff)
	}

	resources, err := mcpClient.ListResources(
		context.Background(),
		mcp.ListResourcesRequest{}, //nolint:exhaustruct
	)
	if err != nil {
		t.Fatalf("failed to list resources: %v", err)
	}

	if diff := cmp.Diff(
		resources,
		//nolint:exhaustruct
		&mcp.ListResourcesResult{
			Resources: []mcp.Resource{},
		},
	); diff != "" {
		t.Errorf("ListResourcesResult mismatch (-want +got):\n%s", diff)
	}

	prompts, err := mcpClient.ListPrompts(
		context.Background(),
		mcp.ListPromptsRequest{}, //nolint:exhaustruct
	)
	if err != nil {
		t.Fatalf("failed to list prompts: %v", err)
	}

	if diff := cmp.Diff(
		prompts,
		//nolint:exhaustruct
		&mcp.ListPromptsResult{
			Prompts: []mcp.Prompt{},
		},
	); diff != "" {
		t.Errorf("ListPromptsResult mismatch (-want +got):\n%s", diff)
	}
}
