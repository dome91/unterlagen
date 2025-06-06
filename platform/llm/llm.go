package llm

import (
	"log/slog"
	"unterlagen/features/assistant"
	"unterlagen/platform/configuration"
)

func GetEmbedder(c configuration.Configuration) assistant.Embedder {
	if c.Assistant.Provider == configuration.OpenAI {
		slog.Info("OpenAI chosen as embedder")
		return NewOpenAI(c)
	}

	if c.Assistant.Provider == configuration.Ollama {
		slog.Info("Ollama chosen as embedder")
		return NewOllama()
	}

	slog.Info("DumbAI chosen as embedder")
	return NewDumbAI()
}

func GetAnswerer(c configuration.Configuration) assistant.Answerer {
	if c.Assistant.Provider == configuration.OpenAI {
		slog.Info("OpenAI chosen as Answerer")
		return NewOpenAI(c)
	}

	if c.Assistant.Provider == configuration.Ollama {
		slog.Info("Ollama chosen as Answerer")
		return NewOllama()
	}

	slog.Info("DumbAI chosen as Answerer")
	return NewDumbAI()
}

func GetChunker(c configuration.Configuration) assistant.Chunker {
	if c.Assistant.Chunker.Type == configuration.Recursive {
		slog.Info("Recursive chosen as Chunker")
		return NewRecursiveChunker(c)
	}

	slog.Info("Fixed Size chosen as Chunker")
	return NewFixedSizeChunker(c)
}
