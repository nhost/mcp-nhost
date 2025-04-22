package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nhost/mcp-nhost/config"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v3"
)

const (
	flagConfigFile = "config-file"
	flagConfirm    = "confirm"
)

func Command() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:  "config",
		Usage: "Writes a sample config file",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:    flagConfigFile,
				Usage:   "Path to the config file",
				Value:   config.GetConfigPath(),
				Sources: cli.EnvVars("CONFIG_FILE"),
			},
			&cli.BoolFlag{ //nolint:exhaustruct
				Name:    flagConfirm,
				Usage:   "Confirm writing the config file",
				Value:   false,
				Sources: cli.EnvVars("CONFIRM"),
			},
		},
		Action: action,
	}
}

//nolint:forbidigo,lll
func action(_ context.Context, cmd *cli.Command) error {
	cfg, err := config.RunWizard()
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to run wizard: %s", err), 1)
	}

	tomlData, err := toml.Marshal(cfg)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to marshal config: %s", err), 1)
	}

	fmt.Println("Generated config file:")
	fmt.Println(string(tomlData))
	fmt.Println("")

	filePath := cmd.String(flagConfigFile)
	fmt.Printf("Now I will write this configuration to the file %s\n", filePath)
	fmt.Println("Proceed? (y/n)")

	var confirm string
	if _, err := fmt.Scanln(&confirm); err != nil {
		return cli.Exit(fmt.Sprintf("failed to read input: %s", err), 1)
	}

	if confirm != "y" && confirm != "Y" {
		fmt.Println("Aborting...")
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil { //nolint:mnd
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0o600); err != nil { //nolint:mnd
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Println("Done! I hope you enjoy using this tool")
	fmt.Println("Please, note that the wizard wasn't exhaustive and you might want to review the documentation to see all the options available. Specially if you want to be more granular on the accesses given to LLMs.")
	return nil
}
