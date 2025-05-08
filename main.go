package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	model "github.com/tokenAPIguy/go-llama/Model"
)

const url = "http://localhost:11434/api/chat"
const chatModel = "qwen2.5-coder"

func submitPrompt(r *model.Request, ch chan<- model.Result) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		ch <- model.Result{Body: nil, Err: err}
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		ch <- model.Result{Body: nil, Err: err}
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- model.Result{Body: nil, Err: err}
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	ch <- model.Result{Body: body, Err: err}
}

func generateRequest(input string) *model.Request {
	return &model.Request{
		Model: chatModel,
		Messages: []model.Message{
			{Role: "user", Content: input},
		},
		Stream: false,
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

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
		block := strings.Join(lines, "\n")

		req := generateRequest(block)
		ch := make(chan model.Result)

		go submitPrompt(req, ch)

		res := <-ch
		if res.Err != nil {
			fmt.Println("Error: ", res.Err)
		}

		var parsed model.Response
		err := json.Unmarshal(res.Body, &parsed)
		if err != nil {
			fmt.Println("Failed to parse response:", err)
			continue
		}
		fmt.Println("<< ", parsed.Message.Content)
	}
}
