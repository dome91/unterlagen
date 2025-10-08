# Unterlagen - Document Archive & Assistant

Unterlagen is a self-hosted document archive system with AI-powered assistant capabilities. Store, organize, and interact with your documents through a clean, modern web interface.

## Features

- **Document Management**: Upload, organize, and search PDF documents with folder structure
- **AI Assistant**: Chat with your documents using OpenAI or Ollama for intelligent document Q&A
- **Document Summarization**: Automatically generate summaries of your documents
- **Export Functionality**: Bulk export of all documents with metadata
- **User Administration**: Secure session-based authentication with user management
- **Modern Interface**: Clean, responsive web interface built with Tailwind CSS and DaisyUI

## Getting Started

### Installation

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

### Docker

Run using Docker:

```bash
# Pull and run the latest image
docker run -d \
  -p 8080:8080 \
  -e UNTERLAGEN_SERVER_SESSION_KEY=your-secret-session-key \
  -v /path/to/your/data:/root/data \
  --name unterlagen \
  ghcr.io/dome91/unterlagen:latest
```

Or with docker-compose:

```yaml
services:
  unterlagen:
    image: ghcr.io/dome91/unterlagen:latest
    ports:
      - "8080:8080"
    environment:
      - UNTERLAGEN_SERVER_SESSION_KEY=your-secret-session-key
    volumes:
      - /path/to/your/data:/root/data
    restart: unless-stopped
```

### Configuration

Configure Unterlagen using environment variables:

**Server Settings:**
- `UNTERLAGEN_SERVER_PORT` - Server port (default: `8080`)
- `UNTERLAGEN_SERVER_BASEURL` - Base URL (default: `http://localhost:8080`)
- `UNTERLAGEN_SERVER_SESSION_KEY` - Session encryption key (**required**)

**AI Assistant Settings:**
- `UNTERLAGEN_ASSISTANT_PROVIDER` - LLM provider: `none`, `openai`, or `ollama` (default: `none`)
- `UNTERLAGEN_ASSISTANT_API_KEY` - API key for OpenAI (required when using OpenAI)

**Ollama-Specific Settings:**
- `UNTERLAGEN_ASSISTANT_OLLAMA_EMBEDDING_MODEL` - Embedding model (default: `nomic-embed-text:latest`)
- `UNTERLAGEN_ASSISTANT_OLLAMA_KNOWLEDGE_BASE_MODEL` - Chat model (default: `phi4:latest`)
- `UNTERLAGEN_ASSISTANT_OLLAMA_SUMMARIZATION_MODEL` - Summarization model (default: `phi4:latest`)

**Example with AI enabled:**
```bash
export UNTERLAGEN_SERVER_SESSION_KEY=your-secret-session-key
export UNTERLAGEN_ASSISTANT_PROVIDER=ollama
export UNTERLAGEN_ASSISTANT_OLLAMA_EMBEDDING_MODEL=nomic-embed-text:latest
export UNTERLAGEN_ASSISTANT_OLLAMA_KNOWLEDGE_BASE_MODEL=phi4:latest
export UNTERLAGEN_ASSISTANT_OLLAMA_SUMMARIZATION_MODEL=phi4:latest
./unterlagen
```

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

## Contributing

Interested in contributing to Unterlagen? See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, architecture details, and contribution guidelines.

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
