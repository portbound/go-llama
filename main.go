package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/tokenAPIguy/go-llama/api"
)

// Temporary env variables
const context = "This is a test environment"
const url = "http://localhost:11434/api/chat"

func saveChat(chat *api.Chat) {
	var f string
	if chat.Name == "" {
		year, month, day := time.Now().Date()
		f = fmt.Sprintf("chats/%d_%d_%d_%s.json", year, month, day, chat.Model)
	}

	jsonData, err := json.Marshal(chat.Messages)
	if err != nil {
		fmt.Println("error marshaling JSON", err)
		return
	}
	err = os.WriteFile(f, jsonData, 0644)
	if err != nil {
		fmt.Printf("Failed to write to %s.", f)
	}
}

func main() {
	// Setting up env
	chatModels := make(map[string]string)
	chatModels["Gemma3"] = "gemma3:1b"
	chatModels["DeepSeek"] = "deepseek-r1:1.5b"
	chat := &api.Chat{}

	client := &api.Client{BaseURL: url}
	var lines string

	fileNames := []huh.Option[string]{}
	entries, err := os.ReadDir("chats")
	if err != nil {
		fmt.Println(err)
	}

	if len(entries) != 0 {
		for _, entry := range entries {
			if !entry.IsDir() { // checking to make sure we don't list subdirs
				name := entry.Name()
				fileNames = append(fileNames, huh.NewOption(name, name))
			}
		}
		resumeChatForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Chat").
					Description("Pick up where you left off?").
					Options(fileNames...).
					Value(&chat.Name),
			),
		)
		err = resumeChatForm.Run()
		if err != nil {
			fmt.Println(err)
		}

		data, err := os.ReadFile("chats/" + chat.Name)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		err = json.Unmarshal(data, &chat.Messages)
		if err != nil {
			fmt.Println("Error parsing JSON into messages:", err)
			return
		}
	}

	if chat.Name == "" {
		opts := []huh.Option[string]{}
		for label, val := range chatModels {
			opts = append(opts, huh.NewOption(label, val))
		}

		modelSelectForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Model").
					Description("Select a model for the chat").
					Options(opts...).
					Value(&chat.Model),

				huh.NewInput().
					Title("Context").
					Description("Provide any additional context for the chat.").
					Value(&chat.Context.Content),
			),
		)
		err := modelSelectForm.Run()
		if err != nil {
			fmt.Println(err)
		}
	}

	// Starting REPL
	for {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Prompt").
					CharLimit(1000).
					Value(&lines),
			),
		)
		err := form.Run()
		if err != nil {
			fmt.Println(err)
		}

		// Build user msg
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "user",
			Content: chat.Context.Content,
		})

		chat.Messages = append(chat.Messages, api.Message{
			Role:    "user",
			Content: lines,
		})

		// Send
		chat.Stream = true

		reply, err := client.HandleRequest(chat)
		if err != nil {
			fmt.Println("Failed to handle client request:", err)
			return
		}

		// Build assistant msg
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "assistant",
			Content: reply.String(),
		})

		saveChat(chat)

		// Reset prompt window
		lines = ""
	}
}
