package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type EmailNotifier struct {
	endpoint      string
	accessKey     string
	senderAddress string
	httpClient    *http.Client
}

// NewEmailNotifier creates a new instance using ACS connection string and sender
func NewEmailNotifier(connectionString, senderAddress string) (*EmailNotifier, error) {
	endpoint, accessKey, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return &EmailNotifier{
		endpoint:      endpoint,
		accessKey:     accessKey,
		senderAddress: senderAddress,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// SendEmailNotification sends a simple email via the ACS REST API
func (n *EmailNotifier) SendEmailNotification(to, subject, plainText, html string) error {
	url := fmt.Sprintf("%s/emails:send?api-version=2023-03-31", n.endpoint)

	payload := map[string]interface{}{
		"senderAddress": n.senderAddress,
		"recipients": map[string]interface{}{
			"to": []map[string]string{
				{"address": to},
			},
		},
		"content": map[string]string{
			"subject":   subject,
			"plainText": plainText,
			"html":      html,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", n.accessKey))

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("email send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("email send failed: status %s", resp.Status)
	}

	return nil
}

// parseConnectionString extracts endpoint and access key from ACS connection string
func parseConnectionString(cs string) (endpoint, accessKey string, err error) {
	parts := strings.Split(cs, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "endpoint=") {
			endpoint = strings.TrimPrefix(part, "endpoint=")
			endpoint = strings.TrimSuffix(endpoint, "/") // clean trailing slash
		}
		if strings.HasPrefix(part, "accesskey=") {
			accessKey = strings.TrimPrefix(part, "accesskey=")
		}
	}
	if endpoint == "" || accessKey == "" {
		err = fmt.Errorf("invalid connection string: must contain endpoint and accesskey")
	}
	return
}
