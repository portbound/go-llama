package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/tokenAPIguy/go-llama/api"
	"github.com/tokenAPIguy/go-llama/forms"
	"github.com/tokenAPIguy/go-llama/ui"
)

const url = "http://localhost:11434/api/chat"

func main() {
	chat, err := initializeChat(&api.Chat{Stream: false})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start REPL
	title := fmt.Sprintf("Prompt - %s", chat.Name)
	var lines string

	// Display chats from history
	if len(chat.Messages) > 0 {
		for _, msg := range chat.Messages {
			fmt.Printf(">> %s:\n", msg.Role)
			fmt.Printf("%s\n\n", msg.Content)
		}
	}

	for {
		// Accept user input
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Placeholder("/bye - to quit").
					Title(title).
					CharLimit(10000).
					Value(&lines),
			),
		)
		err = form.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Quit chat?
		if lines == "/bye" {
			err := chat.SaveChat()
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Build user msg
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "user",
			Content: lines,
		})
		fmt.Println(">> user:")
		fmt.Printf("%s\n\n", chat.Messages[len(chat.Messages)-1].Content)

		// POST to server
		done := make(chan struct{})
		var resp *http.Response

		go func() {
			client := &api.Client{BaseURL: url}
			resp, err = client.HandleRequest(chat)
			if err != nil {
				fmt.Println(err)
				return
			}
			close(done)
		}()

		ui.RunSpinnerUntil(done)

		// Read assistant msg
		assistantMsg, err := readResponse(resp)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Build assistantMsg
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "assistant",
			Content: assistantMsg,
		})
		fmt.Println(">> assistant:")
		fmt.Printf("%s\n\n", assistantMsg)

		// Reset prompt window
		lines = ""
	}
}

func initializeChat(chat *api.Chat) (*api.Chat, error) {
	// Check for existing chats
	entries, err := CheckDir("chats")
	if err != nil {
		return nil, err
	}

	if len(entries) > 0 {
		var files []os.DirEntry
		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, entry)
			}
		}

		// offer to resume an existing chat
		if len(files) != 0 {
			err := forms.ResumeChat(chat, files)
			if err != nil {
				return nil, err
			}
		}
	}

	// if we didn't find any existing chats, or the user declined, start new chat
	if chat.Name == "" {
		err := forms.NewChat(chat)
		if err != nil {
			fmt.Println(err)
		}
	}

	return chat, nil
}

func readResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	var result struct {
		Message api.Message `json:"message"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Message.Content, nil
}
