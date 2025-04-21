# mcp-nhost

A Model Context Protocol (MCP) server implementation for interacting with Nhost Cloud projects and services.

## Overview

MCP-Nhost is designed to provide a unified interface for managing Nhost projects through the Model Context Protocol. It enables seamless interaction with Nhost Cloud services, offering a robust set of tools for project management and configuration.

## Available Tools

The following tools are currently exposed through the MCP interface:

1. **cloud_get_graphql_schema**
   - Provides the GraphQL schema for the Nhost Cloud
   - Provides information about queries, mutations, and type definitions

2. **cloud_graphql_query**
   - Executes GraphQL queries and mutations against the Nhost Cloud
   - Enables project and organization management
   - Allows querying and updating project's configuration

3. **local_get_graphql_schema**
   - Retrieves the GraphQL schema for local Nhost development projects
   - Provides access to project-specific queries and mutations
   - Helps understand available operations for local development helping generating code

4. **local_graphql_query**
   - Executes GraphQL queries against local Nhost development projects
   - Enables testing and development of project-specific operations
   - Supports both queries and mutations for local development

## Screenshots and Examples

You can find screenshots and examples of the current features and tools in the [screenshots](docs/screenshots.md) file.


## Installing

To install mcp-nhost, you can use the following command:

```bash
sudo curl -L https://raw.githubusercontent.com/nhost/mcp-nhost/main/get.sh | bash
```

## Upgrading

To upgrade mcp-nhost, you can use the following command:

```bash
sudo mcp-nhost upgrade --confirm
```

## Getting Started

After installing mcp-nhost, you will need to do two things to get things up and running:

1. Create a Personal Access Token (PAT) in your Nhost account.
2. Configure your mcp-nhost client (i.e. Cursor, etc.)

### Create a PAT

In order to use mcp-nhost, you need to create a Personal Access Token (PAT) in your Nhost account. You can do this by following these steps:

1. Go to https://app.nhost.io/account
2. Scroll don to "Personal Access Tokens" and create a new token. Write down the token as you will need it later.
3. Proceed to configure your mcp-nhost client (instructions for some below).

### Configuring clients

The examples below enable mutations against the cloud and your projects via the flags `--with-cloud-mutations` and `--with-project-mutations`. If you want to disable these, you can remove the flags from the snippets below.

#### Cursor

1. Go to "Cursor Settings"
2. Click on "MCP"
3. Click on "+ Add new global MCP server"
4. Add the following object inside `"mcpServers"`:

```json
    "mcp-nhost": {
      "command": "/usr/local/bin/mcp-nhost",
      "args": [
        "start",
        "--with-cloud-mutations"
      ],
      "env": {
        "NHOST_PAT": "<here-goes-your-pat>"
      }
    }
```

## Roadmap

- ✅ Cloud platform: Basic project and organization management
- ✅ Cloud projects: Configuration management
- ✅ Local projects: Configuration management
- ✅ Local projects: Graphql Schema awareness and query execution
- ✅ Cloud projects: Schema awareness and query execution
- 🔄 Local projects: Create migrations
- 🔄 Local projects: Manage permissions
- 🔄 Documentation: integrate or document use of mintlify's mcp server
- 🔄 Local projects: Auth and Storage schema awareness (maybe via mintlify?)
- 🔄 Cloud projects: Auth and Storage schema awareness (maybe via mintlify?)

If you have any suggestions or feature requests, please feel free to open an issue for discussion.

## Contributing

We welcome contributions to mcp-nhost! If you have suggestions, bug reports, or feature requests, please open an issue or submit a pull request.
