package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tokenAPIguy/go-llama/internal/api"
)

// Temporary env variables
const chatModel = "qwen2.5-coder"
const context = "This is a test environment. Please keep responses concise."
const url = "http://localhost:11434/api/chat"

func main() {
	// Setting up env
	reader := bufio.NewReader(os.Stdin)
	req := api.NewChatRequest(chatModel, context)
	client := &api.Client{BaseURL: url}

	// Starting REPL
	for {
		fmt.Print(">> ")
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

		// Build chat msg
		prompt := strings.Join(lines, "\n")
		newMsg := api.Message{Role: "user", Content: prompt}
		req.Messages = append(req.Messages, newMsg)

		// Send
		res, err := client.HandleRequest(req)
		if err != nil {
			fmt.Println("Failed to handle client request:", err)
			return
		}

		// Receive
		fmt.Println("<< ", res.Message.Content)
	}
}
