package forms

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/tokenAPIguy/go-llama/api"
)

func NewChat(chat *api.Chat) error {
	chatModels := make(map[string]string)
	chatModels["qwen2.5-coder:latest"] = "qwen2.5-coder"

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
				Title("Title").
				Description("Optional: Title of the chat.\nWill default to yyyy_mm_dd_model").
				Value(&chat.Name),
		),
	)
	err := modelSelectForm.Run()
	if err != nil {
		return err
	}

	return nil
}

func ResumeChat(chat *api.Chat, entries []os.DirEntry) error {
	fileNames := []huh.Option[string]{}
	for _, entry := range entries {
		if !entry.IsDir() { // checking to make sure we don't list subdirs
			name := entry.Name()
			fileNames = append(fileNames, huh.NewOption(name, name))
			fileNames = append(fileNames, huh.NewOption("", "New Chat"))
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
	err := resumeChatForm.Run()
	if err != nil {
		return err
	}

	data, err := os.ReadFile("chats/" + chat.Name)
	if err != nil {
		return fmt.Errorf("Failed to read chat log: %w", err)
	}

	err = json.Unmarshal(data, &chat.Messages)
	if err != nil {
		return fmt.Errorf("Failed to parse JSON: %w", err)
	}

	return nil
}
