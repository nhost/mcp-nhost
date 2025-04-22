package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nhost/mcp-nhost/cmd/config"
	"github.com/nhost/mcp-nhost/cmd/gen"
	"github.com/nhost/mcp-nhost/cmd/start"
	"github.com/nhost/mcp-nhost/cmd/upgrade"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var Version = "dev"

func markdownDocs() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:  "docs",
		Usage: "Generate markdown documentation for the CLI",
		Action: func(_ context.Context, cmd *cli.Command) error {
			md, err := docs.ToMarkdown(cmd.Root())
			if err != nil {
				return cli.Exit("failed to generate markdown documentation: "+err.Error(), 1)
			}

			fmt.Println(md) //nolint:forbidigo

			return nil
		},
	}
}

func Command() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:  "nhost-mcp",
		Usage: "Nhost's Model Context Protocol (MCP) server",
		Commands: []*cli.Command{
			markdownDocs(),
			config.Command(),
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
