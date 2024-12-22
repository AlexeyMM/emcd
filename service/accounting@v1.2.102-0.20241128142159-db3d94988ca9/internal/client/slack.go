package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Notifier interface {
	SendMessage(channelName []string, message string)
}

type SlackNotifier struct{}

func (*SlackNotifier) SendMessage(channelNames []string, message string) error {
	ctx := context.Background()
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	jsonMsg := fmt.Sprintf(`{  "text":%q, "username": "LimitBlockBot", "icon_emoji": ":bangbang:" }`, message)
	payload := strings.NewReader(jsonMsg)

	for _, channel := range channelNames {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, channel, payload)

		if err != nil {
			return fmt.Errorf("slack new request: %w", err)
		}

		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("slack request do: %w", err)
		}

		if err = res.Body.Close(); err != nil {
			return fmt.Errorf("slack close: %w", err)
		}
	}

	return nil
}
