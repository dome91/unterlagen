# Unterlagen - Document Archive & Assistant

Unterlagen is a Go web application that serves as a document archive system with integrated AI assistant capabilities. It provides a clean, modern interface for managing documents and leverages LLM integration for intelligent document interaction.

---

## For Users

### What is Unterlagen?

Unterlagen is designed to:
- Provide secure document storage and organization

### Features

- **Document Management**: Upload, organize, and search PDF documents with folder structure
- **User Administration**: Secure session-based authentication and user management
- **Modern Interface**: Clean, responsive web interface built with Tailwind CSS

### Getting Started

#### Installation

1. [Download latest release](https://github.com/Dome91/unterlagen/releases)
2. Set the required session key:
   ```bash
   export UNTERLAGEN_SERVER_SESSION_KEY=your-secret-session-key
   ```
3. Run the application:
   ```bash
   ./unterlagen
   ```
4. Open your browser to `http://localhost:8080`

#### Docker

Run using Docker:

```bash
# Pull and run the latest image
docker run -d \
  -p 8080:8080 \
  -e UNTERLAGEN_SERVER_SESSION_KEY=your-secret-session-key \
  -v /path/to/your/archives:/root/archives \
  --name unterlagen \
  ghcr.io/dome91/unterlagen:latest
```

Or with docker-compose:

```yaml
version: '3.8'
services:
  unterlagen:
    image: ghcr.io/dome91/unterlagen:latest
    ports:
      - "8080:8080"
    environment:
      - UNTERLAGEN_SERVER_SESSION_KEY=your-secret-session-key
      - UNTERLAGEN_ARCHIVE_DIRECTORY=/archives
    volumes:
      - ./archives:/root/archives
    restart: unless-stopped
```

#### Configuration

Configure using environment variables:

- `UNTERLAGEN_SERVER_PORT` - Server port (default: 8080)
- `UNTERLAGEN_SERVER_BASEURL` - Base URL (default: http://localhost:8080)
- `UNTERLAGEN_SERVER_SESSION_KEY` - Session encryption key (required)
- `UNTERLAGEN_ARCHIVE_DIRECTORY` - Document storage path (default: ./archives)

### Roadmap

#### Planned Features
- **Advanced Search**: Full-text search across document content with relevance scoring
- **Document Versioning**: Track changes and maintain document history
- **Collaborative Features**: Multi-user document sharing and commenting
- **OCR Integration**: Extract text from scanned documents and images
- **Bulk Operations**: Mass document operations (tag, move, delete)
- **API Endpoints**: RESTful API for programmatic access
- **Integration Connectors**: Connect to cloud storage (Dropbox, Google Drive, OneDrive)
- **Advanced Analytics**: Document usage patterns and insights dashboard
- **Workflow Automation**: Rule-based document processing and routing

#### AI/ML Enhancements
- **Document Classification**: Automatic categorization using ML models
- **Smart Tagging**: AI-powered tag suggestions
- **Content Summarization**: Automatic document summaries
- **Duplicate Detection**: Find and merge similar documents
- **Sentiment Analysis**: Analyze document tone and sentiment
- **Language Translation**: Multi-language support with translation
- **Voice Integration**: Voice commands and document dictation
- **Predictive Filing**: Suggest optimal folder structure based on content

#### User Experience
- **Dark Mode**: Full dark theme support
- **Drag & Drop**: Enhanced file upload experience
- **Keyboard Shortcuts**: Power user keyboard navigation
- **Customizable Dashboard**: User-configurable widgets and layouts
- **Advanced Filters**: Complex search and filtering capabilities
- **Batch Upload**: Multi-file upload with progress tracking
- **Preview System**: In-browser preview for various file types
- **Notifications**: Real-time alerts and email notifications

---

## For Developers

### Architecture

Unterlagen follows a clean architecture pattern with clear separation of concerns:

```
├── cmd/
│   └── unterlagen.go         # Application entry point & dependency injection
├── features/                 # Business logic organized by domain
│   ├── administration/       # User management, settings, system admin
│   ├── archive/             # Document storage, folders, PDF processing
│   ├── assistant/           # AI chat and document assistance
│   └── common/              # Shared domain logic (ID generation, scheduling)
├── platform/                # Infrastructure layer
│   ├── configuration/       # Viper-based config management
│   ├── database/           # SQLite with goose migrations
│   │   ├── memory/         # In-memory implementations for testing
│   │   └── sqlite/         # Production SQLite implementation
│   ├── event/              # Event system (synchronous, extensible to NATS)
│   ├── llm/                # LLM providers (OpenAI, Ollama)
│   ├── messaging/          # Internal message passing
│   ├── storage/            # File system document storage
│   └── web/                # HTTP server, templates, static assets
├── test/                   # End-to-end tests using Playwright
├── testdata/               # Mock PDFs and test fixtures
└── scripts/                # PDF generation and utility scripts
```

### Technical Features

- **Clean Architecture**: Domain-driven design with dependency injection
- **Type-Safe Templates**: Templ-based HTML templates with Tailwind CSS styling
- **SQLite Database**: Lightweight database with migration support using goose
- **Event System**: Synchronous event handling (extensible to NATS)
- **Testing**: End-to-end tests with Playwright

### Development Setup

#### Prerequisites
- Go 1.21 or later
- Node.js and npm (for Tailwind CSS)

#### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd unterlagen
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Install frontend dependencies:
   ```bash
   cd platform/web
   npm install
   cd ../..
   ```

4. Generate templates and build CSS:
   ```bash
   go generate ./...
   ```

5. Set required environment variables:
   ```bash
   export UNTERLAGEN_SERVER_SESSION_KEY=your-secret-session-key
   ```

6. Run the application:
   ```bash
   # Development with hot reload
   air
   
   # Or build and run manually
   go build -o ./tmp/main cmd/unterlagen.go
   ./tmp/main
   ```

### Dependencies

Key dependencies include:
- **Chi Router**: HTTP routing and middleware
- **Templ**: Type-safe HTML templates
- **Tailwind CSS**: Utility-first CSS framework
- **SQLite**: Embedded database with goose migrations
- **Gorilla Sessions**: Session management
- **Viper**: Configuration management
- **OpenAI/Ollama**: LLM integration

## License

MIT License

Copyright (c) 2025

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
