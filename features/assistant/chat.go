package assistant

type Role string

const (
	UserRole      Role = "USER"
	AssistantRole Role = "ASSISTANT"
)

type ChatMessage struct {
	Role    Role
	Content string
}

type Chat struct {
	ID       string
	Messages []ChatMessage
	UserID   string
}

func (c *Chat) AddMessage(role Role, message string) {
	c.Messages = append(c.Messages, ChatMessage{
		Role:    role,
		Content: message,
	})
}

type ChatRepository interface {
	Save(chat Chat) error
	FindByIDAndUserID(id string, userID string) (Chat, error)
}
