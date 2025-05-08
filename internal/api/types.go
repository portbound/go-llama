package api

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChatRequest(model string, context string) *ChatRequest {
	return &ChatRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: context,
			},
		},
		Stream: false,
	}
}
