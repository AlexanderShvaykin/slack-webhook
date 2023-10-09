package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type slackMessage struct {
	Text        string       `json:"text,omitempty"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
}

type SendMessage interface {
  PostMessage(url string, text string, mrkdwn bool) ([]byte, error)
}

type Client struct { }

func (c *Client) PostMessage(url string, text string, mrkdwn bool) ([]byte, error) {
  msg := slackMessage {
    Text: text,
    Markdown: mrkdwn,
  }

  marshalled, err := json.Marshal(msg)
  if err != nil {
    return nil, err
  }

  resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalled))
  if err != nil {
    return nil, err
  }

  body, err := io.ReadAll(resp.Body) 
  if err != nil {
    return nil, err
  }

  if resp.StatusCode > 300 {
    return nil, errors.New(string(body))
  }

  return body, nil
}

