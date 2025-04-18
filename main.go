package main

import (
	"context"
	"log"
	"os"

	"github.com/nhost/mcp-nhost/cmd/gen"
	"github.com/nhost/mcp-nhost/cmd/mcp"
	"github.com/urfave/cli/v3"
)

var Version = "dev"

func main() {
	cmd := &cli.Command{ //nolint:exhaustruct
		Name:  "nhost-mcp",
		Usage: "Nhost's Model Context Protocol (MCP) server",
		Commands: []*cli.Command{
			mcp.Command(),
			gen.Command(),
		},
		Version: Version,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
