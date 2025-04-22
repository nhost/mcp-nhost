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
   - Uses "user" role unless specified otherwise

4. **local_graphql_query**
   - Executes GraphQL queries against local Nhost development projects
   - Enables testing and development of project-specific operations
   - Supports both queries and mutations for local development
   - Uses "user" role unless specified otherwise

5. **local_config_server_schema**
   - Retrieves the GraphQL schema for the local config server
   - Helps understand available configuration options

6. **local_config_server_query**
   - Executes GraphQL queries against the local config server
   - Enables querying and modifying local project configuration
   - Changes require running 'nhost up' to take effect

7. **project_get_graphql_schema**
   - Retrieves the GraphQL schema for Nhost Cloud projects
   - Provides access to project-specific queries and mutations
   - Uses "user" role unless specified otherwise

8. **project_graphql_query**
   - Executes GraphQL queries against Nhost Cloud projects
   - Enables interaction with live project data
   - Supports both queries and mutations (need to be allowed)
   - Uses "user" role unless specified otherwise

## Screenshots and Examples

You can find screenshots and examples of the current features and tools in the [screenshots](docs/screenshots.md) file.

## Installing

To install mcp-nhost, you can use the following command:

```bash
sudo curl -L https://raw.githubusercontent.com/nhost/mcp-nhost/main/get.sh | bash
```

## Configuring

After installing mcp-nhost, you will need to configure it. You can do this by running the command `mcp-nhost config` in your terminal. See [CONFIG.md](docs/CONFIG.md) for more details.

## Configuring clients

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
      ],
    }
```

## CLI Usage

For help on how to use the CLI, you can run:

```bash
mcp-nhost --help
```

Or check [USAGE.md](docs/USAGE.md) for more details.

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
