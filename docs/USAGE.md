# NAME

nhost-mcp - Nhost's Model Context Protocol (MCP) server

# SYNOPSIS

nhost-mcp

```
[--help|-h]
[--version|-v]
```

**Usage**:

```
nhost-mcp [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--help, -h**: show help

**--version, -v**: print the version


# COMMANDS

## docs

Generate markdown documentation for the CLI

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## start

Starts the MCP server

**--bind**="": Bind address in the form <host>:<port>. If omitted use stdio

**--help, -h**: show help

**--local-admin-secret**="": Admin secret for local projects (default: nhost-admin-secret)

**--local-config-server-url**="": Config server URL for local projects (default: https://local.dashboard.local.nhost.run/v1/configserver/graphql)

**--local-graphql-url**="": GraphQL URL for local projects (default: https://local.graphql.local.nhost.run/v1)

**--nhost-pat**="": Personal Access Token

**--project-admin-secret**="": Admin secret for the project

**--project-allow-mutations**="": Allow mutations for the project. If empty, no mutations are allowed. Use * to allow all mutations. Can be passed multiple times (default: [])

**--project-allow-queries**="": Allow queries for the project. If empty, no queries are allowed. Use * to allow all queries. Can be passed multiple times (default: [])

**--project-pat**="": Personal Access Token for the project

**--project-region**="": Region for the project

**--project-subdomain**="": Subdomain for the project

**--with-cloud-mutations**: Enable mutations against Nhost Cloud to allow operating on projects

### help, h

Shows a list of commands or help for one command

## gen

Generate GraphQL schema for Nhost Cloud

**--help, -h**: show help

**--nhost-pat**="": Personal Access Token

**--with-mutations**: Include mutations in the generated schema

### help, h

Shows a list of commands or help for one command

## upgrade

Checks if there is a new version and upgrades it

**--confirm**: Confirm the upgrade without prompting

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## help, h

Shows a list of commands or help for one command

