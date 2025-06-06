package llm

import (
	"context"
	"strings"
	"unterlagen/features/assistant"

	"github.com/ollama/ollama/api"
)

var (
	_ assistant.Embedder = &Ollama{}
	_ assistant.Answerer = &Ollama{}
)

const (
	embeddingModel     = "nomic-embed-text"
	knowledgeBaseModel = "gemma:2b"
)

type Ollama struct {
	client *api.Client
}

// Generate implements assistant.Embedder.
func (o *Ollama) Generate(text string) (assistant.Embeddings, error) {
	response, err := o.client.Embeddings(context.Background(), &api.EmbeddingRequest{
		Model:  embeddingModel,
		Prompt: text,
	})
	if err != nil {
		return nil, err
	}

	return response.Embedding, nil
}

// Answer implements assistant.Answerer.
func (o *Ollama) Answer(question string, nodes []assistant.Node) (string, error) {
	var contextForQuestion string
	for _, node := range nodes {
		contextForQuestion += node.Chunk
		contextForQuestion += "\n"
	}
	systemMessage := strings.ReplaceAll(prompt, "{context}", contextForQuestion)

	var answer string
	err := o.client.Generate(context.Background(), &api.GenerateRequest{
		Model:  knowledgeBaseModel,
		System: systemMessage,
		Prompt: question,
		Stream: new(bool),
	}, func(gr api.GenerateResponse) error {
		answer = gr.Response
		return nil
	})

	return answer, err
}

func NewOllama() *Ollama {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}

	return &Ollama{client: client}
}
