package llm

import (
	"context"
	"encoding/json"
	"strings"
	"unterlagen/features/archive"
	"unterlagen/features/assistant"
	"unterlagen/platform/configuration"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/packages/param"
	"github.com/openai/openai-go/v2/shared"
)

var (
	_ assistant.Embedder         = &OpenAI{}
	_ assistant.Answerer         = &OpenAI{}
	_ archive.DocumentSummarizer = &OpenAI{}
)

const assistantPrompt = `
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

const summarizePrompt = `
You are a document analyzer. Analyze the document and provide a structured summary.

Guidelines:
- Create a one-sentence overview describing what the document is about
- List 3-5 key points, facts, or topics covered in the document
- Use objective, professional language
- Focus on the most important content and conclusions
- Do not add information not present in the original text

Document text:
{text}
`

var dimension int64 = 768

var DocumentSummaryResponseSchema = GenerateSchema[archive.DocumentSummary]()

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

	systemMessage := strings.ReplaceAll(assistantPrompt, "{context}", contextForQuestion)
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

func (o *OpenAI) SummarizeText(text string) (archive.DocumentSummary, error) {
	systemMessage := strings.ReplaceAll(summarizePrompt, "{text}", text)
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemMessage),
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "document_summary",
		Description: openai.String("Summary of the document with overview and key points"),
		Schema:      DocumentSummaryResponseSchema,
		Strict:      openai.Bool(true),
	}

	response, err := o.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       shared.ChatModelGPT4oMini, // Use mini for cost efficiency
		Temperature: param.NewOpt(0.3),         // Slightly more creative for better summaries
		MaxTokens:   param.NewOpt[int64](500),  // Limit summary length
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})

	if err != nil {
		return archive.DocumentSummary{}, err
	}

	var summary archive.DocumentSummary
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &summary)
	if err != nil {
		return archive.DocumentSummary{}, err
	}

	return summary, nil
}

func NewOpenAI(configuration configuration.Configuration) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(option.WithAPIKey(configuration.Assistant.ApiKey)),
	}
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}
