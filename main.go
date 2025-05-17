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
	chat := &api.Chat{Stream: false}
	client := &api.Client{BaseURL: url}

	// Check for existing chats
	entries, err := CheckDir("chats")

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
			fmt.Println(err)
			return
		}
	}

	// if we didn't find any existing chats, or the user declined, start new chat
	if chat.Name == "" {
		err := forms.NewChat(chat)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Start REPL
	title := fmt.Sprintf("Prompt - %s", chat.Name)
	var lines string

	if len(chat.Messages) > 0 {
		for _, msg := range chat.Messages {
			fmt.Printf(">> %s:\n", msg.Role)
			fmt.Printf("%s\n\n", msg.Content)
		}
	}

	for {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
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

		// Build user msg
		fmt.Println(chat.Messages)
		chat.Messages = append(chat.Messages, api.Message{
			Role:    "user",
			Content: lines,
		})

		fmt.Println(chat.Messages)
		// Send
		done := make(chan struct{})
		var resp *http.Response

		go func() {
			resp, err = client.HandleRequest(chat)
			close(done)
		}()

		ui.RunSpinnerUntil(done)

		fmt.Println(">> user:")
		fmt.Printf("%s\n\n", chat.Messages[len(chat.Messages)-1].Content)

		fullReply, err := renderResponse(resp)
		if err != nil {
			fmt.Println(err)
			return
		}

		chat.Messages = append(chat.Messages, api.Message{
			Role:    "assistant",
			Content: fullReply,
		})

		fmt.Println(">> assistant:")
		fmt.Printf("%s\n\n", fullReply)

		// TODO: figure out why it's duplicating existing files with .json, WriteFile should create only if it doesnt exist...

		err = chat.SaveChat()
		if err != nil {
			fmt.Println(err)
		}

		// Reset prompt window
		lines = ""
	}
}

func renderResponse(resp *http.Response) (string, error) {
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
