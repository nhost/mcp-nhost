package main

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/nhost/mcp-nhost/cmd/start"
	"github.com/nhost/mcp-nhost/tools/cloud"
)

func TestStart(t *testing.T) { //nolint:cyclop
	t.Parallel()

	cmd := Command()

	buf := bytes.NewBuffer(nil)
	cmd.Writer = buf

	go func() {
		if err := cmd.Run(
			context.Background(),
			[]string{"main", "start", "--bind=:9000", "--nhost-pat=asdasd"},
		); err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)

	transportClient, err := transport.NewSSEClientTransport("http://localhost:9000/sse")
	if err != nil {
		t.Fatalf("failed to create transport client: %v", err)
	}

	mcpClient, err := client.NewClient(transportClient)
	if err != nil {
		t.Fatalf("failed to create mcp client: %v", err)
	}
	defer mcpClient.Close()

	if diff := cmp.Diff(
		mcpClient.GetServerInfo(),
		protocol.Implementation{
			Name:    "nhost-mcp",
			Version: "dev",
		},
	); diff != "" {
		t.Errorf("ServerInfo mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(
		mcpClient.GetServerCapabilities(),
		protocol.ServerCapabilities{
			Prompts: &protocol.PromptsCapability{
				ListChanged: true,
			},
			Resources: &protocol.ResourcesCapability{
				ListChanged: true,
				Subscribe:   true,
			},
			Tools: &protocol.ToolsCapability{
				ListChanged: true,
			},
		},
	); diff != "" {
		t.Errorf("ServerCapabilities mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(
		mcpClient.GetServerInstructions(),
		start.ServerInstructions,
	); diff != "" {
		t.Errorf("ServerInstructions mismatch (-want +got):\n%s", diff)
	}

	tools, err := mcpClient.ListTools(context.Background())
	if err != nil {
		t.Fatalf("failed to list tools: %v", err)
	}
	if diff := cmp.Diff(
		tools,
		&protocol.ListToolsResult{
			Tools: []*protocol.Tool{
				{
					Name:        "cloud-get-graphql-schema",
					Description: cloud.ToolGetGraphqlSchemaInstructions,
					InputSchema: protocol.InputSchema{
						Type:       "object",
						Properties: nil,
						Required:   nil,
					},
				},
				{
					Name:        "cloud-graphql-query",
					Description: cloud.ToolGraphqlQueryInstructions,
					InputSchema: protocol.InputSchema{
						Type: "object",
						Properties: map[string]*protocol.Property{
							"query": {
								Type:        "string",
								Description: "graphql query to perform",
							},
							"variables": {
								Type:        "string",
								Description: "variables to use in the query",
							},
						},
						Required: []string{"query"},
					},
				},
			},
			NextCursor: "",
		},
		cmpopts.SortSlices(func(a, b *protocol.Tool) bool {
			return a.Name < b.Name
		}),
	); diff != "" {
		t.Errorf("ListToolsResult mismatch (-want +got):\n%s", diff)
	}

	resources, err := mcpClient.ListResources(context.Background())
	if err != nil {
		t.Fatalf("failed to list resources: %v", err)
	}

	if diff := cmp.Diff(
		resources,
		&protocol.ListResourcesResult{
			Resources:  []protocol.Resource{},
			NextCursor: "",
		},
	); diff != "" {
		t.Errorf("ListResourcesResult mismatch (-want +got):\n%s", diff)
	}
}
