package llm

import (
	"context"
	"encoding/json"
	"strings"
	"unterlagen/features/archive"
	"unterlagen/features/assistant"
	"unterlagen/platform/configuration"

	"github.com/ollama/ollama/api"
)

var (
	_ assistant.Embedder         = &Ollama{}
	_ assistant.Answerer         = &Ollama{}
	_ archive.DocumentSummarizer = &Ollama{}
)

type Ollama struct {
	client              *api.Client
	summarizationFormat json.RawMessage
	embeddingModel      string
	knowledgeBaseModel  string
	summarizationModel  string
}

// Generate implements assistant.Embedder.
func (o *Ollama) Generate(text string) (assistant.Embeddings, error) {
	response, err := o.client.Embeddings(context.Background(), &api.EmbeddingRequest{
		Model:  o.embeddingModel,
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
	systemMessage := strings.ReplaceAll(assistantPrompt, "{context}", contextForQuestion)

	var answer string
	err := o.client.Generate(context.Background(), &api.GenerateRequest{
		Model:  o.knowledgeBaseModel,
		System: systemMessage,
		Prompt: question,
		Stream: new(bool),
	}, func(gr api.GenerateResponse) error {
		answer = gr.Response
		return nil
	})

	return answer, err
}

func (o *Ollama) SummarizeText(text string) (archive.DocumentSummary, error) {
	systemPrompt := `You are a document analyzer. Analyze the document and provide a structured summary.
					Provide an overview that is one clear, concise sentence describing what the document is about.
					Provide key_points that list 3-5 important facts, topics, or conclusions from the document.
					Guidelines:
					- Use objective, professional language
					- Focus on the most important content and conclusions
					- Do not add information not present in the original text`

	var response string
	err := o.client.Generate(context.Background(), &api.GenerateRequest{
		Model:  o.summarizationModel,
		System: systemPrompt,
		Prompt: text,
		Stream: new(bool),
		Format: o.summarizationFormat,
	}, func(gr api.GenerateResponse) error {
		response = gr.Response
		return nil
	})

	if err != nil {
		return archive.DocumentSummary{}, err
	}
	var summary archive.DocumentSummary
	err = json.Unmarshal([]byte(response), &summary)
	return summary, err
}

func NewOllama(config configuration.Configuration) *Ollama {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}

	return &Ollama{
		client:              client,
		summarizationFormat: []byte(`
	{
	  "type": "object",
	  "properties": {
		"overview": {
		  "type": "string"
		},
		"key_points": {
		  "type": "array",
		  "items": {
			"type": "string"
		  }
		}
	  }
	}`),
		embeddingModel:      config.Assistant.Ollama.EmbeddingModel,
		knowledgeBaseModel:  config.Assistant.Ollama.KnowledgeBaseModel,
		summarizationModel:  config.Assistant.Ollama.SummarizationModel,
	}
}
