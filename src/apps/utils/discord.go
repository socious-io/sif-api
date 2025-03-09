package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func DiscordSendTextMessage(webhookURL, message string) error {

	payload, err := json.Marshal(map[string]string{
		"content": message,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send log to Discord, status: %s", resp.Status)
	}

	return nil
}
