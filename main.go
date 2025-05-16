package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/tokenAPIguy/go-llama/api"
	"github.com/tokenAPIguy/go-llama/forms"
)

const url = "http://localhost:11434/api/chat"

func main() {
	chat := &api.Chat{}
	client := &api.Client{BaseURL: url}

	// Check for existing chats
	entries, err := os.ReadDir("chats")
	//TODO: should create dir if it doesn't exist
	if err != nil {
		fmt.Println(err)
	}

	// If there are existing chats, offer to resume one of them
	if len(entries) != 0 {
		err := forms.ResumeChat(chat, entries)
		if err != nil {
			fmt.Println(err)
		}
	}

	// If an existing chat was not resumed, create a new one
	if chat.Name == "" {
		err := forms.NewChat(chat)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Starting REPL
	var lines string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Prompt").
				CharLimit(1000).
				Value(&lines),
		),
	)
	err = form.Run()
	if err != nil {
		fmt.Println(err)
	}

	for {
		// Build user msg
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "user",
			Content: lines,
		})

		// Send
		chat.Stream = true

		reply, err := client.HandleRequest(chat)
		if err != nil {
			fmt.Println("Failed to handle client request:", err)
		}

		// Build assistant msg
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "assistant",
			Content: reply.String(),
		})

		err = chat.SaveChat()
		if err != nil {
			fmt.Println(err)
		}

		// Reset prompt window
		lines = ""
	}
}

//TODO: dynamically pull the available models
