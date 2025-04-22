package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

func ptr[T any](v T) *T {
	return &v
}

const (
	DefaultLocalConfigServerURL = "https://local.dashboard.local.nhost.run/v1/configserver/graphql"
	DefaultLocalGraphqlURL      = "https://local.graphql.local.nhost.run/v1"
)

type Config struct {
	// If configured allows managing the cloud. For instance, this allows you to configure
	// projects, list projects, organizations, and so on.
	Cloud *Cloud `json:"cloud,omitempty" toml:"cloud"`

	// If configured allows working with a local project running via the CLI. This includes
	// configuring it, working with the schema, migrations, etc.
	Local *Local `json:"local,omitempty" toml:"local"`

	// Projects is a list of projects that you want to allow access to. This grants access to the
	// GraphQL schema allowing it to inspect it and run allowed queries and mutations.
	Projects []Project `json:"projects" toml:"projects"`
}

type Cloud struct {
	// Personal Access Token to authenticate with the Nhost Cloud API. You can get one
	// on the following URL: https://app.nhost.io/account
	PAT string `json:"pat" toml:"pat"`

	// If enabled you can run mutations against the Nhost Cloud to manipulate project's configurations
	// amongst other things. Queries are always allowed if this section is configured.
	EnableMutations bool `json:"enable_mutations" toml:"enable_mutations"`
}

type Local struct {
	// Admin secret to use when running against a local project.
	AdminSecret string `json:"admin_secret" toml:"admin_secret"`

	// GraphQL URL to use when running against a local project.
	// Defaults to "https://local.dashboard.local.nhost.run/v1/configserver/graphql"
	ConfigServerURL *string `json:"config_server_url,omitempty" toml:"config_server_url,omitempty"`

	// GraphQL URL to use when running against a local project.
	// Defaults to "https://local.graphql.local.nhost.run/v1"
	GraphqlURL *string `json:"graphql_url,omitempty" toml:"graphql_url,omitempty"`
}

type Project struct {
	// Project's subdomain
	Subdomain string `json:"subdomain" toml:"subdomain"`

	// Project's region
	Region string `json:"region" toml:"region"`

	// Admin secret to operate against the project.
	// Either admin secret or PAT is required.
	AdminSecret *string `json:"admin_secret,omitempty" toml:"admin_secret,omitempty"`

	// PAT to operate against the project. Note this PAT must belong to this project.
	// Either admin secret or PAT is required.
	PAT *string `json:"pat,omitempty" toml:"pat,omitempty"`

	// List of queries that are allowed to be executed against the project.
	// If empty, no queries are allowed. Use [*] to allow all queries.
	AllowQueries []string `json:"allow_queries" toml:"allow_queries"`

	// List of mutations that are allowed to be executed against the project.
	// If empty, no mutations are allowed. Use [*] to allow all mutations.
	// Note that this is only used if the project is configured to allow mutations.
	AllowMutations []string `json:"allow_mutations" toml:"allow_mutations"`
}

func GetConfigPath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "mcp-nhost.toml"
		}
		configHome = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configHome, "nhost", "mcp-nhost.toml")
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	if config.Local != nil {
		if config.Local.GraphqlURL == nil {
			config.Local.GraphqlURL = ptr(DefaultLocalGraphqlURL)
		}
		if config.Local.ConfigServerURL == nil {
			config.Local.ConfigServerURL = ptr(DefaultLocalConfigServerURL)
		}
	}

	return &config, nil
}
