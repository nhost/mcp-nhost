package gen

import (
	"context"
	"fmt"

	"github.com/nhost/mcp-nhost/graphql"
	"github.com/nhost/mcp-nhost/nhost/auth"
	"github.com/urfave/cli/v3"
)

const (
	flagNhostAuthURL    = "nhost-auth-url"
	flagNhostGraphqlURL = "nhost-graphql-url"
	flagNhostPAT        = "nhost-pat"
)

func Command() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:  "gen",
		Usage: "Generate GraphQL schema for Nhost Cloud",
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
		},
		Action: action,
	}
}

func action(ctx context.Context, cmd *cli.Command) error {
	interceptor, err := auth.WithPAT(
		cmd.String(flagNhostAuthURL),
		cmd.String(flagNhostPAT),
	)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	var introspection graphql.ResponseIntrospection
	if err := graphql.Query(
		ctx,
		cmd.String(flagNhostGraphqlURL),
		graphql.IntrospectionQuery,
		nil,
		&introspection,
		interceptor,
	); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	schema := graphql.ParseSchema(
		introspection, graphql.Filter{
			AllowQueries: []graphql.Queries{
				{
					Name:           "organizations",
					DisableNesting: true,
				},
				{
					Name:           "organization",
					DisableNesting: true,
				},
				{
					Name:           "app",
					DisableNesting: true,
				},
				{
					Name:           "apps",
					DisableNesting: true,
				},
				{
					Name:           "config",
					DisableNesting: false,
				},
			},
			AllowMutations: []graphql.Queries{
				{
					Name:           "updateConfig",
					DisableNesting: false,
				},
			},
		},
	)

	fmt.Print(schema) //nolint:forbidigo

	return nil
}
