package assistant

import (
	"unterlagen/features/archive"
	"unterlagen/features/common"
)

type Answerer interface {
	Answer(question string, nodes []Node) (string, error)
}

type Embeddings []float64
type Embedder interface {
	Generate(text string) (Embeddings, error)
}

type Chunker interface {
	Chunk(text string) ([]string, error)
}

type Node struct {
	ID         string
	Chunk      string
	Embeddings Embeddings
	DocumentID string
}

type NodeRepository interface {
	SaveAll(nodes []Node) error
	FindSimiliarByEmbedding(embeddings Embeddings) ([]Node, error)
	DeleteAllByDocumentID(documentID string) error
}

type Assistant struct {
	nodeRepository NodeRepository
	chatRepository ChatRepository
	answerer       Answerer
	embedder       Embedder
	chunker        Chunker
}

func (a *Assistant) StartChat(userID string) (Chat, error) {
	chat := Chat{
		ID:       common.GenerateID(),
		Messages: []ChatMessage{},
		UserID:   userID,
	}

	err := a.chatRepository.Save(chat)
	return chat, err
}

func (a *Assistant) GetChat(chatID string, userID string) (Chat, error) {
	return a.chatRepository.FindByIDAndUserID(chatID, userID)
}

func (a *Assistant) Ask(chatID string, question string, userID string) error {
	chat, err := a.chatRepository.FindByIDAndUserID(chatID, userID)
	if err != nil {
		return err
	}

	embedding, err := a.embedder.Generate(question)
	if err != nil {
		return err
	}

	nodes, err := a.nodeRepository.FindSimiliarByEmbedding(embedding)
	if err != nil {
		return err
	}

	answer, err := a.answerer.Answer(question, nodes)
	if err != nil {
		return err
	}

	chat.AddMessage(UserRole, question)
	chat.AddMessage(AssistantRole, answer)
	return a.chatRepository.Save(chat)
}

func (a *Assistant) generateNodes(document archive.Document) error {
	chunks, err := a.chunker.Chunk(document.Text)
	if err != nil {
		return err
	}

	var nodes []Node
	for _, chunk := range chunks {
		embeddings, err := a.embedder.Generate(chunk)
		if err != nil {
			return err
		}

		nodes = append(nodes, Node{
			ID:         common.GenerateID(),
			Chunk:      chunk,
			Embeddings: embeddings,
			DocumentID: document.ID,
		})
	}

	return a.nodeRepository.SaveAll(nodes)
}

func (a *Assistant) deleteNodes(document archive.Document) error {
	return a.nodeRepository.DeleteAllByDocumentID(document.ID)
}

func NewAssistant(
	nodeRepository NodeRepository,
	chatRepository ChatRepository,
	answerer Answerer,
	embedder Embedder,
	chunker Chunker,
	documentMessages archive.DocumentMessages,
) *Assistant {
	assistant := &Assistant{
		nodeRepository: nodeRepository,
		chatRepository: chatRepository,
		answerer:       answerer,
		embedder:       embedder,
		chunker:        chunker,
	}

	documentMessages.SubscribeDocumentAnalyzed(assistant.generateNodes)
	documentMessages.SubscribeDocumentDeleted(assistant.deleteNodes)

	return assistant
}
