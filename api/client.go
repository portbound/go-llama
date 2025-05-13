package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL string
}

var fullReply strings.Builder

func (c *Client) HandleRequest(req *Chat) (strings.Builder, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return fullReply, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fullReply, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fullReply, err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	fmt.Printf(">> %s << ", req.Model)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fullReply, err
		}

		// Optional: filter empty lines
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		var chunk struct {
			Message Message `json:"message"`
			Done    bool    `json:"done"`
		}

		err = json.Unmarshal(line, &chunk)
		if err != nil {
			return fullReply, fmt.Errorf("error parsing stream chunk: %w\nChunk: %s", err, line)
		}

		fullReply.WriteString(chunk.Message.Content)
		fmt.Print(chunk.Message.Content)

		if chunk.Done {
			break
		}
	}

	fmt.Println()
	fmt.Println()
	return fullReply, nil
}
