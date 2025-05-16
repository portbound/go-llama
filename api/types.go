package api

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Chat struct {
	Name     string    `json:"name"`
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

// func NewChat(model string, context string) *Chat {
// 	return &Chat{
// 		Model: model,
// 		Messages: []Message{
// 			{
// 				Role:    "user",
// 				Content: context,
// 			},
// 		},
// 		Stream: true,
// 	}
// }

func (c *Chat) SaveChat() error {
	if c.Name == "" {
		year, month, day := time.Now().Date()
		c.Name = fmt.Sprintf("%d_%d_%d_%s", year, month, day, c.Model)
	}

	jsonData, err := json.Marshal(c.Messages)
	if err != nil {
		return err
	}

	f := fmt.Sprintf("chats/%s.json", c.Name)
	err = os.WriteFile(f, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
