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
# Build Tailwind CSS (from platform/web directory) and JavaScript
cd platform/web
npm run build:css
npm run build:js
```

### Testing
```bash
# Run tests (requires session key environment variable)
UNTERLAGEN_SERVER_SESSION_KEY=my-key go test ./test/...
```

### Database Migrations
Migrations are in `platform/database/sqlite/migrations/` and use goose for version management. Migrations are automatically applied during application startup.

## Configuration

Configuration uses Viper with environment variables prefixed with `UNTERLAGEN_`:
- `UNTERLAGEN_SERVER_PORT` (default: 8080)
- `UNTERLAGEN_SERVER_BASEURL` (default: http://localhost:8080)
- `UNTERLAGEN_SERVER_SESSION_KEY` (required)
- `UNTERLAGEN_ASSISTANT_PROVIDER` (none, openai, ollama)
- `UNTERLAGEN_ASSISTANT_API_KEY`
- `UNTERLAGEN_ASSISTANT_OLLAMA_EMBEDDING_MODEL` (default: embeddinggemma:300m)
- `UNTERLAGEN_ASSISTANT_OLLAMA_KNOWLEDGE_BASE_MODEL` (default: phi4:latest)
- `UNTERLAGEN_ASSISTANT_OLLAMA_SUMMARIZATION_MODEL` (default: phi4:latest)

## Template System

Uses templ for type-safe HTML templates. Run `templ generate` or `go generate ./...` to compile templates before building.
Use HTMX for client-side and partial updates.

## Testing Strategy

Uses Playwright for end-to-end testing with in-memory sqlite repositories for isolation. Tests require the session key environment variable to be set.

## Frontend Patterns and Best Practices

### Template Structure
- **templ Components**: Use semantic component names (e.g., `DocumentActions`, `DocumentTitleDisplay`)
- **ID Naming**: Use kebab-case for HTML IDs (e.g., `title-section`, `action-buttons`)
- **CSS Classes**: Combine Tailwind utility classes with DaisyUI component classes
- **Responsive Design**: Use responsive prefixes (`lg:`, `md:`) for layout adjustments

### JavaScript Integration
- **Inline Scripts**: Place JavaScript directly in templ `script` blocks for component-specific functionality
- **Global Functions**: Expose functions to `window` object when needed across components
- **DOM Manipulation**: Use `innerHTML` for dynamic content replacement
- **Event Handling**: Attach event listeners via `onclick` attributes or `addEventListener`

### Form Handling
- **POST + Redirect Pattern**: Use standard form POST with server redirect instead of complex HTMX partial updates for simplicity
- **Hidden Inputs**: Use hidden form fields to pass data between JavaScript and server
- **Form Validation**: Leverage HTML5 validation attributes (`required`, etc.)

### UI State Management
- **Mode-Based Interfaces**: Toggle between view/edit modes by replacing HTML content
- **Action Button Transformation**: Replace action button sets based on current mode (e.g., Edit/Download/Trash ↔ Save/Cancel)
- **Focus Management**: Auto-focus and select text in inputs when entering edit mode
- **Escape Key Handling**: Provide keyboard shortcuts for canceling operations

### Document Management Patterns
- **Title Editing**: Implement inline editing with mode switching rather than separate edit pages
- **File Operations**: Use direct links for downloads, forms for state-changing operations
- **Preview Integration**: Embed document previews with navigation controls
- **Status Indicators**: Use visual cues (opacity, icons) to show document states (trashed, etc.)

### Testing Patterns
- **ARIA Role Selectors**: ALWAYS prefer `GetByRole()` with ARIA roles (e.g., `role="alert"`, `role="button"`) over CSS selectors for better accessibility and maintainability
- **Page Object Pattern**: Use Playwright locators and role-based selectors
- **Test Flow Integration**: Combine multiple operations (upload → edit → download) in single test flows
- **State Verification**: Verify both UI state changes and data persistence
- **Wait Strategies**: Use `WaitForLoadState` for network-dependent operations

### Error Handling and UX
- **Confirmation Dialogs**: Use browser `confirm()` for destructive operations
- **Progressive Enhancement**: Ensure forms work without JavaScript
- **Loading States**: Handle page transitions with appropriate loading indicators
- **User Feedback**: Provide clear success/error messages via flash notifications
