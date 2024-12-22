package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"code.emcdtech.com/emcd/sdk/log"
)

type Slack interface {
	Send(ctx context.Context, text string) error
}

type slack struct {
	cli        *http.Client
	webhookURL string
}

func NewSlack(webhookURL string) *slack {
	return &slack{
		cli:        &http.Client{},
		webhookURL: webhookURL,
	}
}

func (s *slack) Send(ctx context.Context, text string) error {
	message := map[string]string{"text": text}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(messageBytes))
	if err != nil {
		return fmt.Errorf("newRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.cli.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error(ctx, "slack: close conn: %s", err.Error())
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Warn(ctx, "slack: status code: %d, message: %s", resp.StatusCode, text)
	}

	return nil
}
