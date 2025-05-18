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
	//TODO: When we reopen chats, we haven't actually saved the model anywhere so subseqeunt requests will not have a model attacahed to them. Try seeing if we can save chats as chats/modelname/chatname.json instead of just chats/chatname.json -- then we can try to pull this out of the path when we do ResumeChat() in main.go
	if c.Name == "" {
		year, month, day := time.Now().Date()
		c.Name = fmt.Sprintf("%d_%d_%d_%s", year, month, day, c.Model)
	}

	jsonData, err := json.Marshal(c.Messages)
	if err != nil {
		return err
	}

	f := fmt.Sprintf("chats/%s/%s.json", c.Model, c.Name)
	err = os.WriteFile(f, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
