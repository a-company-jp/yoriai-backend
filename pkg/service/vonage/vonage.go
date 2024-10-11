package vonage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/a-company/yoriai-backend/pkg/config"
	"io"
	"log/slog"
	"net/http"
)

const (
	API_URL = "https://studio-api-us.ai.vonage.com/telephony/make-call"
)

type VonageService struct {
	conf config.VonageConfig
}

func NewVonage() *VonageService {
	return &VonageService{
		conf: config.Config.VonageConfig,
	}
}

type PhoneAPIInput struct {
	PhoneNumber   string
	ReceiverName  string
	CallerName    string
	RemindMessage string
}

func (s VonageService) CallPhoneAPI(input PhoneAPIInput) error {
	slog.Info("CallPhoneAPI", slog.Any("input", input))
	requestBody, err := createRequestBody(s.conf.AgentID, input)
	if err != nil {
		return err
	}

	responseBody, err := sendApiRequest(requestBody, s.conf.VgaiKey)
	if err != nil {
		return err
	}

	slog.Info("Response Body:", slog.String("response", responseBody))
	return nil
}

func createRequestBody(agentId string, input PhoneAPIInput) ([]byte, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"agent_id": agentId,
		"to":       input.PhoneNumber,
		"session_parameters": []map[string]string{
			{
				"name":  "RECEIVER_NAME",
				"value": input.ReceiverName,
			},
			{
				"name":  "CALLER_NAME",
				"value": input.CallerName,
			},
			{
				"name":  "REMIND_MESSAGE",
				"value": input.RemindMessage,
			},
		},
	})

	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}
	return requestBody, nil
}

func sendApiRequest(requestBody []byte, XKey string) (string, error) {
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		slog.Error("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vgai-Key", XKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", err)
		return "", err
	}

	return string(body), nil
}
