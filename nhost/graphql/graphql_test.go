package graphql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nhost/mcp-nhost/nhost/auth"
	"github.com/nhost/mcp-nhost/nhost/graphql"
)

func ptr[T any](v T) *T {
	return &v
}

func TestGraphql(t *testing.T) {
	cl, err := graphql.NewClientWithResponses(
		"https://local.hasura.local.nhost.run",
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var up graphql.MigrationStep
	if err := up.FromSqlMigrationStep(graphql.SqlMigrationStep{
		Type: graphql.RunSql,
		Args: graphql.SqlMigrationArgs{
			Cascade:  ptr(true),
			ReadOnly: ptr(false),
			Source:   "default",
			Sql:      "CREATE TABLE IF NOT EXISTS test (id serial PRIMARY KEY, name text);",
		},
	}); err != nil {
		t.Fatalf("failed to create up migrations: %s", err)
	}

	var down graphql.MigrationStep
	if err := down.FromSqlMigrationStep(graphql.SqlMigrationStep{
		Type: graphql.RunSql,
		Args: graphql.SqlMigrationArgs{
			Cascade:  ptr(true),
			ReadOnly: ptr(false),
			Source:   "default",
			Sql:      "DROP TABLE IF EXISTS test;",
		},
	}); err != nil {
		t.Fatalf("failed to create down migrations: %s", err)
	}

	resp, err := cl.ExecuteMigrationWithResponse(
		context.Background(),
		graphql.ExecuteMigrationJSONRequestBody{
			Datasource: "default",
			Name:       "create_test_table",
			Up:         []graphql.MigrationStep{up},
			Down:       &[]graphql.MigrationStep{down},
		},
		auth.WithAdminSecret("nhost-admin-secret"),
	)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode() != 200 {
		t.Fatalf(
			"unexpected status code: %d\n%s, %s",
			resp.StatusCode(),
			*resp.JSON500.Code,
			*resp.JSON500.Message,
		)
	}

	fmt.Println("Response:", *resp.JSON200.Name)

	var track graphql.MigrationStep
	if err := track.FromTrackTableStep(graphql.TrackTableStep{
		Type: graphql.PgTrackTable,
		Args: graphql.TrackTableArgs{
			Source: "default",
			Table: struct {
				Name   string "json:\"name\""
				Schema string "json:\"schema\""
			}{
				Name:   "test",
				Schema: "public",
			},
		},
	}); err != nil {
		t.Fatalf("failed to create track migrations: %s", err)
	}

	resp, err = cl.ExecuteMigrationWithResponse(
		context.Background(),
		graphql.ExecuteMigrationJSONRequestBody{
			Datasource: "default",
			Name:       "track_test_table",
			Up:         []graphql.MigrationStep{track},
			Down:       &[]graphql.MigrationStep{},
		},
	)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected status code: %d", resp.StatusCode())
	}

	fmt.Println("Response:", *resp.JSON200.Name)
}
