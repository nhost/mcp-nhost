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

## Getting Started

TODO

## Contributing

TODO:w
