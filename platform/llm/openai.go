package llm

import (
	"context"
	"strings"
	"unterlagen/features/assistant"
	"unterlagen/platform/configuration"

	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	_ assistant.Embedder = &OpenAI{}
	_ assistant.Answerer = &OpenAI{}
)

const prompt = `
You are a helpful assistant with access to specific document context. Your role is to:

1. Answer questions based ONLY on the provided context
2. If the context doesn't contain enough information to answer fully, acknowledge this clearly
3. Quote relevant parts of the context to support your answers
4. Never make up information beyond what's in the context
5. If you're unsure about something, say so directly

Context:
{context}

Please provide a clear, accurate answer based solely on the context above. Include relevant quotes when appropriate.
`

var dimension int64 = 768

type OpenAI struct {
	client openai.Client
}

func (o *OpenAI) Generate(text string) (assistant.Embeddings, error) {
	response, err := o.client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(text),
		},
		Model:          openai.EmbeddingModelTextEmbedding3Small,
		Dimensions:     param.NewOpt(dimension),
		EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
	})

	if err != nil {
		return nil, err
	}

	return response.Data[0].Embedding, nil
}

func (o *OpenAI) Answer(question string, nodes []assistant.Node) (string, error) {
	var contextForQuestion string
	for _, node := range nodes {
		contextForQuestion += node.Chunk
		contextForQuestion += "\n"
	}

	systemMessage := strings.ReplaceAll(prompt, "{context}", contextForQuestion)
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemMessage),
		openai.UserMessage(question),
	}
	response, err := o.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       shared.ChatModelGPT4o,
		Temperature: param.NewOpt(0.),
	})

	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

func NewOpenAI(configuration configuration.Configuration) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(option.WithAPIKey(configuration.Assistant.ApiKey)),
	}
}
