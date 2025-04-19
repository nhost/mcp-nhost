package main

import (
	"context"
	"log"
	"os"

	"github.com/nhost/mcp-nhost/cmd/gen"
	"github.com/nhost/mcp-nhost/cmd/start"
	"github.com/nhost/mcp-nhost/cmd/upgrade"
	"github.com/urfave/cli/v3"
)

var Version = "dev"

func Command() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:  "nhost-mcp",
		Usage: "Nhost's Model Context Protocol (MCP) server",
		Commands: []*cli.Command{
			start.Command(),
			gen.Command(),
			upgrade.Command(),
		},
		Version: Version,
	}
}

func main() {
	cmd := Command()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
