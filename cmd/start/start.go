package start

import (
	"context"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/pkg"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/nhost/mcp-nhost/nhost/auth"
	"github.com/nhost/mcp-nhost/tools/cloud"
	"github.com/urfave/cli/v3"
)

const (
	flagNhostAuthURL       = "nhost-auth-url"
	flagNhostGraphqlURL    = "nhost-graphql-url"
	flagNhostPAT           = "nhost-pat"
	flagBind               = "bind"
	flagWithCloudMutations = "with-cloud-mutations"
)

const (
	// this seems to be largely ignored by clients, or at least by cursor.
	// we also need to look into roots and resources as those might be helpful.
	ServerInstructions = `
		This is an MCP server to interact with Nhost Cloud and with projects running on it.

		Important notes to anyone using this MCP server. Do not use this MCP server without
		following these instructions:

		1. Before attempting to call any tool *-graphql-query, always get the schema using the
		   *-get-graphql-schema tool
		2. Apps and projects are the same and while users may talk about projects in the GraphQL
		  api those are referred as apps.
		3. IDs are always UUIDs so if you have anything else (like an app/project name) you may need
		   to first get the ID using the *-graphql-query tool.
		4. If you have an error querying the GraphQL API, please check the schema again. The schema may
		   have changed and the query you are using may be invalid.
	`
)

func Command() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:  "start",
		Usage: "Starts the MCP server",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:    flagNhostAuthURL,
				Usage:   "Nhost auth URL",
				Hidden:  true,
				Value:   "https://otsispdzcwxyqzbfntmj.auth.eu-central-1.nhost.run/v1",
				Sources: cli.EnvVars("NHOST_AUTH_URL"),
			},
			&cli.StringFlag{ //nolint:exhaustruct
				Name:    flagNhostGraphqlURL,
				Usage:   "Nhost GraphQL URL",
				Hidden:  true,
				Value:   "https://otsispdzcwxyqzbfntmj.graphql.eu-central-1.nhost.run/v1",
				Sources: cli.EnvVars("NHOST_GRAPHQL_URL"),
			},
			&cli.StringFlag{ //nolint:exhaustruct
				Name:     flagNhostPAT,
				Usage:    "Personal Access Token",
				Required: true,
				Sources:  cli.EnvVars("NHOST_PAT"),
			},
			&cli.StringFlag{ //nolint:exhaustruct
				Name:     flagBind,
				Usage:    "Bind address in the form <host>:<port>. If omitted use stdio",
				Required: false,
				Sources:  cli.EnvVars("BIND"),
			},
			&cli.BoolFlag{ //nolint:exhaustruct
				Name:    flagWithCloudMutations,
				Usage:   "Enable mutations against Nhost Cloud to allow operating on projects",
				Value:   false,
				Sources: cli.EnvVars("WITH_CLOUD_MUTATIONS"),
			},
		},
		Action: action,
	}
}

func getLogger(debug bool) pkg.Logger { //nolint:ireturn
	var logger pkg.Logger
	if debug {
		logger = pkg.DebugLogger
	} else {
		logger = pkg.DefaultLogger
	}

	return logger
}

func action(_ context.Context, cmd *cli.Command) error {
	logger := getLogger(cmd.Bool("debug"))

	var transportServer transport.ServerTransport
	var err error
	if bind := cmd.String(flagBind); bind != "" {
		logger.Infof("listening on " + bind)
		transportServer, err = transport.NewSSEServerTransport(bind)
	} else {
		logger.Infof("listening on stdio")
		transportServer = transport.NewStdioServerTransport()
	}
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to create transport server: %v", err), 1)
	}

	mcpServer, err := server.NewServer(
		transportServer,
		server.WithServerInfo(
			protocol.Implementation{
				Name:    cmd.Root().Name,
				Version: cmd.Root().Version,
			},
		),
		server.WithLogger(logger),
		server.WithInstructions(ServerInstructions),
	)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to create mcp server: %v", err), 1)
	}

	interceptor, err := auth.WithPAT(
		cmd.String(flagNhostAuthURL),
		cmd.String(flagNhostPAT),
	)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	cloudTool := cloud.NewTool(
		cmd.String(flagNhostGraphqlURL), cmd.Bool(flagWithCloudMutations), interceptor,
	)

	if err := cloudTool.Register(mcpServer); err != nil {
		return cli.Exit(fmt.Sprintf("failed to register Nhost's cloud tools: %v", err), 1)
	}

	logger.Infof("starting mcp server")
	if err = mcpServer.Run(); err != nil {
		return cli.Exit(fmt.Sprintf("failed to run mcp server: %v", err), 1)
	}

	return nil
}
