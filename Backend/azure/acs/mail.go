package acs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NLstn/clubs/azure"
)

type EmailRequest struct {
	SenderAddress string     `json:"senderAddress"`
	Recipients    Recipients `json:"recipients"`
	Content       Content    `json:"content"`
}

type Recipients struct {
	To []Recipient `json:"to"`
}

type Recipient struct {
	Address string `json:"address"`
}

type Content struct {
	Subject   string `json:"subject"`
	PlainText string `json:"plainText"`
	HTML      string `json:"html,omitempty"`
}

func SendMail(senderAddress string, recipients []Recipient, subject, plainText, htmlContent string) error {
	// Get authentication token from auth.go
	token := azure.GetACSToken()
	if token == "" {
		return fmt.Errorf("failed to get ACS token for sending email")
	}

	// Get ACS endpoint from environment variable
	endpoint := os.Getenv("AZURE_ACS_ENDPOINT")
	if endpoint == "" {
		return fmt.Errorf("AZURE_ACS_ENDPOINT environment variable not set")
	}

	// Create email request
	emailReq := EmailRequest{
		SenderAddress: senderAddress,
		Recipients: Recipients{
			To: recipients,
		},
		Content: Content{
			Subject:   subject,
			PlainText: plainText,
			HTML:      htmlContent,
		},
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/emails:send?api-version=2023-03-31", endpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	} else {
		return fmt.Errorf("failed to send email. Status code: %d", resp.StatusCode)
	}
}

func SendTestMail() {
	senderAddress := "DoNotReply@9a1159be-1efb-4f84-b252-9e949a999910.azurecomm.net"
	recipients := []Recipient{
		{"niklas.lahnstein@outlook.com"},
	}
	subject := "Test Email from ACS"
	plainText := "This is a test email sent using Azure Communication Services."
	htmlContent := "<html><body><h1>Test Email</h1><p>This is a test email sent using Azure Communication Services.</p></body></html>"

	err := SendMail(senderAddress, recipients, subject, plainText, htmlContent)
	if err != nil {
		log.Default().Println(err.Error())
		return
	}

	log.Default().Println("Test email sent successfully")
}
