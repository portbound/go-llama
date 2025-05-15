package api

type Chat struct {
	Name     string    `json:"name"`
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Context  Message   `json:"context"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChat(model string, context string) *Chat {
	return &Chat{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: context,
			},
		},
		Stream: true,
	}
}
