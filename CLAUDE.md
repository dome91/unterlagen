# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Unterlagen is a Go web application template that serves as a document archive system with assistant/chat capabilities. It follows a clean architecture pattern with features, platform, and infrastructure layers.

## Architecture

### Core Structure
- **cmd/**: Application entry point
- **features/**: Business logic organized by domain (administration, archive, assistant)
- **platform/**: Infrastructure concerns (database, event, storage, web, configuration)
- **test/**: End-to-end tests using Playwright

### Key Components
- **Database**: SQLite with goose migrations
- **Web Framework**: Chi router with templ templates and Tailwind CSS
- **Event System**: Synchronous events (can be swapped for NATS)
- **Storage**: Filesystem-based document storage
- **Authentication**: Session-based with Gorilla sessions
- **LLM Integration**: Supports OpenAI and Ollama providers

### Dependency Injection Pattern
The main function in `cmd/unterlagen.go` wires up all dependencies manually, following a dependency injection pattern without a DI container.

## Development Commands

### Building and Running
```bash
# Development with hot reload
air

# Manual build and run
go build -o ./tmp/main cmd/unterlagen.go
UNTERLAGEN_SERVER_SESSION_KEY=unterlagen ./tmp/main

# Generate templ templates and CSS
go generate ./...
```

### Frontend Assets
```bash
# Build Tailwind CSS (from platform/web directory)
cd platform/web
npm run build:css
```

### Testing
```bash
# Run tests (requires session key environment variable)
UNTERLAGEN_SERVER_SESSION_KEY=my-key go test ./test/...
```

### Database Migrations
Migrations are in `platform/database/sqlite/migrations/` and use goose for version management.

## Configuration

Configuration uses Viper with environment variables prefixed with `UNTERLAGEN_`:
- `UNTERLAGEN_SERVER_PORT` (default: 8080)
- `UNTERLAGEN_SERVER_BASEURL` (default: http://localhost:8080)
- `UNTERLAGEN_SERVER_SESSION_KEY` (required)
- `UNTERLAGEN_ASSISTANT_PROVIDER` (none, openai, ollama)
- `UNTERLAGEN_ASSISTANT_API_KEY`

## Template System

Uses templ for type-safe HTML templates. Run `templ generate` or `go generate ./...` to compile templates before building.

## Testing Strategy

Uses Playwright for end-to-end testing with in-memory sqlite repositories for isolation. Tests require the session key environment variable to be set.