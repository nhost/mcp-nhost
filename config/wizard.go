package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//nolint:forbidigo
func RunWizard() (*Config, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Nhost MCP Configuration Wizard!")
	fmt.Println("==============================================")
	fmt.Println()

	cloudConfig := wizardCloud(reader)
	fmt.Println("")
	localConfig := wizardLocal(reader)
	fmt.Println("")
	projects := wizardProject(reader)
	fmt.Println("")

	return &Config{
		Cloud:    cloudConfig,
		Local:    localConfig,
		Projects: projects,
	}, nil
}

//nolint:forbidigo,lll
func wizardCloud(reader *bufio.Reader) *Cloud {
	fmt.Println("1. Access to the Nhost Cloud platform allows LLMs to access your Nhost projects and organizations configuration. It is useful to view your projects, configure them and to perform other cloud operations you may normally do via the dashboard.")
	if promptYesNo(reader, "- Do you want to enable access to Nhost Cloud?") {
		pat := promptString(reader, "* Enter your Personal Access Token (PAT) from https://app.nhost.io/account:")
		return &Cloud{
			PAT:             pat,
			EnableMutations: true,
		}
	}

	return nil
}

//nolint:forbidigo,lll
func wizardLocal(reader *bufio.Reader) *Local {
	fmt.Println("2. Access to local development allows LLMs to access your Nhost local development environment allowing them to view and manipulate your project configuration and to access your local GraphQL API. This can give LLMs context to generate code that interacts with your local Nhost project.")

	if promptYesNo(reader, "- Do you want to enable access to local development?") {
		adminSecret := promptString(reader, "* Enter your Admin Secret for local development (typically nhost-admin-secret):")

		return &Local{
			AdminSecret:     adminSecret,
			ConfigServerURL: nil,
			GraphqlURL:      nil,
		}
	}

	return nil
}

//nolint:forbidigo,lll
func wizardProject(reader *bufio.Reader) []Project {
	projects := make([]Project, 0)
	fmt.Println("3. Access to specific projects allows LLMs to access your Nhost projects' GraphQL API. This is useful to query and manipulate your projects' data. You can restrict which queries and mutations are allowed to be performed on each project. Visit the documentation for more information about this.")
	if promptYesNo(reader, "- Do you want to configure access to specific projects?") {
		for {
			project := Project{
				Subdomain:      "",
				Region:         "",
				AdminSecret:    nil,
				PAT:            nil,
				AllowQueries:   []string{"*"},
				AllowMutations: []string{"*"},
			}

			project.Subdomain = promptString(reader, "* Enter project subdomain:")
			project.Region = promptString(reader, "* Enter project region:")

			authType := promptChoice(reader, "Choose authentication method:", []string{"Admin Secret", "PAT"})
			if authType == "Admin Secret" {
				adminSecret := promptString(reader, "Enter project Admin Secret:")
				project.AdminSecret = &adminSecret
			} else {
				pat := promptString(reader, "Enter project PAT:")
				project.PAT = &pat
			}

			projects = append(projects, project)

			if !promptYesNo(reader, "Do you want to add another project?") {
				break
			}
		}
	}
	return projects
}

//nolint:forbidigo
func promptString(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt + " ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

//nolint:forbidigo
func promptYesNo(reader *bufio.Reader, prompt string) bool {
	for {
		fmt.Printf("%s (y/n) ", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		fmt.Println("Please answer with 'y' or 'n'")
	}
}

//nolint:forbidigo
func promptChoice(reader *bufio.Reader, prompt string, options []string) string {
	for {
		fmt.Printf("%s\n", prompt)
		for i, opt := range options {
			fmt.Printf("%d) %s\n", i+1, opt)
		}
		fmt.Print("Enter number: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if num := strings.TrimSpace(input); num != "" {
			switch num {
			case "1":
				return options[0]
			case "2":
				return options[1]
			}
		}
		fmt.Println("Please select a valid option")
	}
}
