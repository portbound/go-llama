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
const cyan = "\033[33m"
const reset = "\033[0m"

func main() {
	chat := &api.Chat{Stream: false}
	err := initializeChat(chat)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start REPL
	title := fmt.Sprintf("Prompt - %s", chat.Name)
	var lines string

	// Display chats from history
	if len(chat.Messages) > 0 {
		for _, msg := range chat.Messages {
			if msg.Role == "assistant" {
				fmt.Printf(cyan+">> %s:\n"+reset, chat.Model)
				fmt.Printf(cyan+"%s\n\n"+reset, msg.Content)
			}
			if msg.Role == "user" {
				fmt.Printf(">> %s:\n", msg.Role)
				fmt.Printf("%s\n\n", msg.Content)
			}
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
		fmt.Printf(cyan+">> %s:\n"+reset, chat.Model)
		fmt.Printf(cyan+"%s\n\n"+reset, assistantMsg)

		// Reset prompt window
		lines = ""
	}
}

func initializeChat(chat *api.Chat) error {
	// Check for existing chats
	entries, err := CheckDir("chats")
	if err != nil {
		return err
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
			err = forms.ResumeChat(chat, files)
			if err != nil {
				return err
			}
		}
	}

	// if we didn't find any existing chats, or the user declined, start new chat
	if chat.Name == "" {
		err := forms.NewChat(chat)
		if err != nil {
			return err
		}
	}

	return nil
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
