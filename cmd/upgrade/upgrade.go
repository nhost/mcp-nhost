package upgrade

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/nhost/mcp-nhost/software"
	"github.com/urfave/cli/v3"
)

const (
	flagConfirm = "confirm"
)

const devVersion = "dev"

func Command() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:    "upgrade",
		Aliases: []string{},
		Usage:   "Checks if there is a new version and upgrades it",
		Action:  action,
		Flags: []cli.Flag{
			&cli.BoolFlag{ //nolint:exhaustruct
				Name:  flagConfirm,
				Usage: "Confirm the upgrade without prompting",
				Value: false,
			},
		},
	}
}

func action(ctx context.Context, cmd *cli.Command) error {
	curVersion := cmd.Root().Version

	mgr := software.NewManager()
	releases, err := mgr.GetReleases(ctx, curVersion)
	if err != nil {
		return fmt.Errorf("failed to get releases: %w", err)
	}

	if len(releases) == 0 {
		fmt.Println("No releases found.") //nolint:forbidigo
		return nil
	}

	latest := releases[0]
	if latest.TagName == curVersion {
		fmt.Println("You are already on the latest version. Hurray!") //nolint:forbidigo
		return nil
	}

	printVersionsSince(releases, curVersion)

	if !cmd.Bool(flagConfirm) {
		fmt.Printf("Do you want to upgrade to %s? (y/N): ", latest.TagName) //nolint:forbidigo
		var answer string
		_, err := fmt.Scanln(&answer)
		if err != nil || answer != "y" && answer != "Y" {
			fmt.Println("Upgrade cancelled.") //nolint:forbidigo
			return nil                        //nolint:nilerr
		}
	}

	return install(ctx, mgr, latest, curVersion)
}

func printVersionsSince(
	releases software.Releases,
	curVersion string,
) {
	fmt.Printf( //nolint:forbidigo
		"Versions released since your current version (%s):\n\n", curVersion,
	)
	for _, release := range releases {
		if release.TagName == curVersion {
			break
		}
		fmt.Printf("# %s\n\n", release.TagName) //nolint:forbidigo
		fmt.Println(release.Body)               //nolint:forbidigo
		fmt.Println()                           //nolint:forbidigo
	}
}

func findAsset(
	release software.Release,
) (string, error) {
	want := fmt.Sprintf("mcp-nhost-%s-%s-%s.tar.gz", release.TagName, runtime.GOOS, runtime.GOARCH)
	for _, asset := range release.Assets {
		if asset.Name == want {
			return asset.BrowserDownloadURL, nil
		}
	}
	return "", fmt.Errorf("failed to find asset for %s", want) //nolint:goerr113
}

func install(
	ctx context.Context,
	mgr *software.Manager,
	latest software.Release,
	curVersion string,
) error {
	fmt.Println("Downloading...") //nolint:forbidigo

	url, err := findAsset(latest)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "mcp-nhost-")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if err := mgr.DownloadAsset(ctx, url, tmpFile); err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	curBin, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find installed mcp-nhost: %w", err)
	}

	if curVersion == devVersion || curVersion == "" {
		curBin = "/tmp/mcp-nhost"
	}

	fmt.Printf("Copying to %s...\n", curBin) //nolint:forbidigo
	if err := os.Rename(tmpFile.Name(), curBin); err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", tmpFile.Name(), curBin, err)
	}

	fmt.Println("Setting permissions...")           //nolint:forbidigo
	if err := os.Chmod(curBin, 0o755); err != nil { //nolint:mnd
		return fmt.Errorf("failed to set permissions on %s: %w", curBin, err)
	}

	fmt.Println("Upgrade complete!") //nolint:forbidigo

	return nil
}
