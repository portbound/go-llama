package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL string
}

var fullReply strings.Builder

func (c *Client) HandleRequest(chat *Chat) (*http.Response, error) {
	jsonBody, err := json.Marshal(chat)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
