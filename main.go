package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tokenAPIguy/go-llama/api"
)

// Temporary env variables
const chatModel = "gemma3:1b"
const context = "This is a test environment"
const url = "http://localhost:11434/api/chat"

func main() {
	// Setting up env
	reader := bufio.NewReader(os.Stdin)
	chat := api.NewChat(chatModel, context)
	client := &api.Client{BaseURL: url}

	// Starting REPL
	for {
		fmt.Print(">> User << ")
		var lines []string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			line = strings.TrimRight(line, "\r\n")
			if line == "" {
				break
			}

			lines = append(lines, line)
		}

		// Build user msg
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "user",
			Content: strings.Join(lines, "\n"),
		})

		// Send
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

		fmt.Printf("Messages in Slice: %d", len(chat.Messages))
	}
}
