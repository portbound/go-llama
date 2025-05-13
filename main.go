package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/tokenAPIguy/go-llama/api"
)

// Temporary env variables
const context = "This is a test environment"
const url = "http://localhost:11434/api/chat"

func main() {
	// Setting up env
	chatModels := make(map[string]string)
	chatModels["Gemma3"] = "gemma3:1b"
	// 	chat := api.NewChat(chatModel, context)
	chat := api.Chat{}
	client := &api.Client{BaseURL: url}
	var lines string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Model").
				Description("Select a model for the chat").
				Options(
					huh.NewOption("Gemma3", chatModels["Gemma3"]),
				).
				Value(&chat.Model),

			huh.NewInput().
				Title("Context").
				Description("Provide any additional context for the chat.").
				Value(&chat.Context.Content),
		),
	)
	err := form.Run()
	if err != nil {
		fmt.Println(err)
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
		reply, err := client.HandleRequest(&chat)
		if err != nil {
			fmt.Println("Failed to handle client request:", err)
			return
		}

		// Build assistant msg
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "assistant",
			Content: reply.String(),
		})

		fmt.Printf("Messages in Slice: %d", len(chat.Messages))
	}

}
