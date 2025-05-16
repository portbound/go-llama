package forms

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/huh"
	"github.com/tokenAPIguy/go-llama/api"
)

func NewChat(chat *api.Chat) error {
	// Lookup installed models
	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(os.Getenv("USERPROFILE"), ".ollama", "models", "manifests", "registry.ollama.ai", "library")
	} else {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".ollama", "models", "manifests", "registry.ollama.ai", "library")
	}

	models, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Append installed models to form select list
	opts := []huh.Option[string]{}
	for _, m := range models {
		//TODO: we're listing subdirs here, so we will want to add support later in case we have more than one subtype of model. Need to check the depth of models/ and if > 1 loop through those to get names
		name := m.Name()
		opts = append(opts, huh.NewOption(name, name))
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
				Description("Optional: If not provided the title will default to yyyy_mm_dd_model").
				Value(&chat.Name),
		),
	)
	err = modelSelectForm.Run()
	if err != nil {
		return err
	}

	return nil
}

func ResumeChat(chat *api.Chat, chats []os.DirEntry) error {
	fileNames := []huh.Option[string]{
		huh.NewOption("New Chat", ""),
	}
	for _, c := range chats {
		if !c.IsDir() { // checking to make sure we only list files
			name := c.Name()
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
	err := resumeChatForm.Run()
	if err != nil {
		return err
	}

	// Early return if someone doesn't actually pick a chat
	if chat.Name == "" {
		return nil
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
