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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Chat) SaveChat() error {
	if c.Name == "" {
		year, month, day := time.Now().Date()
		c.Name = fmt.Sprintf("%d_%d_%d_%s", year, month, day, c.Model)
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("chats/%s", c.Name)
	f := fmt.Sprintf("%s.json", name)
	err = os.WriteFile(f, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
