# mcp-nhost

A Model Control Protocol (MCP) server implementation for interacting with Nhost Cloud services. This project provides a standardized interface for managing and interacting with Nhost projects, both in the cloud and locally.

## Overview

MCP-Nhost is designed to provide a unified interface for managing Nhost projects through the Model Control Protocol. It enables seamless interaction with Nhost Cloud services, offering a robust set of tools for project management and configuration.

## Current Features

### Nhost Cloud

#### Project Management
- Query and manage Nhost projects
- Manage project resources and settings

### Organization Management
- Query organizations

#### Configuration Management
- Access and modify project configurations

## Available Tools

The following tools are currently exposed through the MCP interface:

1. **Nhost GraphQL Schema Access**
   - Retrieve the complete GraphQL schema for Nhost Cloud
   - Access type definitions and available operations

2. **Nhost GraphQL Query Execution**
   - Execute queries against Nhost Cloud
   - Perform operations on projects and organizations
   - Manage configurations and settings

## Roadmap

### Phase 1: Cloud Integration (Current)
- ✅ Basic project and organization management
- ✅ Configuration management

### Phase 2: CLI Integration (Coming Soon)
- 🔄 Local project management
- 🔄 Development environment configuration

### Phase 3: Development Workflows (Planned)
- 🔄 Manage migrations
- 🔄 Manage permissions
- 🔄 Schema-awareness for developmenet workflows

### Phase 4: Production integration (Planned)
- 🔄 Integration with production projects via Graphite for building agents effortlessly


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

## Contributing

We welcome contributions to mcp-nhost! If you have suggestions, bug reports, or feature requests, please open an issue or submit a pull request.
