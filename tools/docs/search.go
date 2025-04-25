package docs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nhost/mcp-nhost/tools"
	"github.com/nhost/mcp-nhost/tools/docs/mintlify"
)

const (
	ToolSearchName = "search"
	//nolint:lll
	ToolSearchInstructions = `Search Nhost's official documentation. Use this tool to look for information about Nhost's features, APIs, guides, etc. Follow relevant links to get more details.`
)

func (t *Tool) registerSearch(mcpServer *server.MCPServer) {
	configServerSchemaTool := mcp.NewTool(
		ToolSearchName,
		mcp.WithDescription(ToolSearchInstructions),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Search Nhost Docs",
				ReadOnlyHint:    true,
				DestructiveHint: false,
				IdempotentHint:  true,
				OpenWorldHint:   true,
			},
		),
		mcp.WithString(
			"query",
			mcp.Description("The search query"),
			mcp.Required(),
		),
	)
	mcpServer.AddTool(configServerSchemaTool, t.handleSearch)
}

func (t *Tool) handleSearch(
	ctx context.Context, req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	query, err := tools.FromParams[string](req.Params.Arguments, "query")
	if err != nil {
		return nil, err
	}

	resp, err := t.mintlify.Autocomplete(
		ctx,
		mintlify.AutocompleteRequest{
			Query:          query,
			PageSize:       10, //nolint:mnd
			SearchType:     "full_text",
			ExtendResults:  true,
			ScoreThreshold: 1,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error calling mintlify: %w", err)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshalling response: %w", err)
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
