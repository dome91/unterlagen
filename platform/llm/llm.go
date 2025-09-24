package llm

import (
	"log/slog"
	"unterlagen/features/archive"
	"unterlagen/features/assistant"
	"unterlagen/platform/configuration"
)

func GetEmbedder(config configuration.Configuration) assistant.Embedder {
	if config.Assistant.Provider == configuration.OpenAI {
		slog.Info("OpenAI chosen as embedder")
		return NewOpenAI(config)
	}

	if config.Assistant.Provider == configuration.Ollama {
		slog.Info("Ollama chosen as embedder")
		return NewOllama(config)
	}

	slog.Info("DumbAI chosen as embedder")
	return NewDumbAI()
}

func GetAnswerer(config configuration.Configuration) assistant.Answerer {
	if config.Assistant.Provider == configuration.OpenAI {
		slog.Info("OpenAI chosen as Answerer")
		return NewOpenAI(config)
	}

	if config.Assistant.Provider == configuration.Ollama {
		slog.Info("Ollama chosen as Answerer")
		return NewOllama(config)
	}

	slog.Info("DumbAI chosen as Answerer")
	return NewDumbAI()
}

func GetChunker(config configuration.Configuration) assistant.Chunker {
	if config.Assistant.Chunker.Type == configuration.Recursive {
		slog.Info("Recursive chosen as Chunker")
		return NewRecursiveChunker(config)
	}

	slog.Info("Fixed Size chosen as Chunker")
	return NewFixedSizeChunker(config)
}

func GetSummarizer(config configuration.Configuration) archive.DocumentSummarizer {
	if config.Assistant.Provider == configuration.OpenAI {
		slog.Info("OpenAI chosen as Summarizer")
		return NewOpenAI(config)
	}

	if config.Assistant.Provider == configuration.Ollama {
		slog.Info("Ollama chosen as Summarizer")
		return NewOllama(config)
	}

	slog.Info("DumbAI chosen as Summarizer")
	return NewDumbAI()
}
