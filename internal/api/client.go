package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
}

func (c *Client) HandleRequest(req *ChatRequest) (ChatResponse, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return ChatResponse{}, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return ChatResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return ChatResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, err
	}

	var parsed ChatResponse
	err = json.Unmarshal(body, &parsed)
	return parsed, err
}
