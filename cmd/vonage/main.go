package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	API_URL = "https://studio-api-us.ai.vonage.com/telephony/make-call"
)

type LINEInput struct {
	phoneNumber   string
	receiverName  string
	callerName    string
	remindMessage sql.NullString
}

func main() {
	input := LINEInput{
		phoneNumber:   "1234567890",
		receiverName:  "Alisa",
		callerName:    "Bob",
		remindMessage: sql.NullString{String: "Don't forget to bring your umbrella", Valid: true},
	}
	err := CallPhoneApi(input)
	if err != nil {
		fmt.Println("Error making API call:", err)
	}
}

func CallPhoneApi(input LINEInput) error {
	XKey, err := FetchEnvVar("XKey")
	if err != nil {
		return err
	}

	AgentId, err := FetchEnvVar("AGENT_ID")
	if err != nil {
		return err
	}

	requestBody, err := CreateRequestBody(AgentId, input)
	if err != nil {
		return err
	}

	responseBody, err := SendApiRequest(requestBody, XKey)
	if err != nil {
		return err
	}

	fmt.Println("Response Body:", responseBody)
	return nil
}

func FetchEnvVar(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return "", err
	}

	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is not set in the .env file", key)
	}
	return value, nil
}

func CreateRequestBody(agentId string, input LINEInput) ([]byte, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"agent_id": agentId,
		"to":       input.phoneNumber,
		"session_parameters": []map[string]string{
			{
				"name":  "RECEIVER_NAME",
				"value": input.receiverName,
			},
			{
				"name":  "CALLER_NAME",
				"value": input.callerName,
			},
			{
				"name":  "REMIND_MESSAGE",
				"value": input.remindMessage.String,
			},
		},
	})

	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}
	return requestBody, nil
}

func SendApiRequest(requestBody []byte, XKey string) (string, error) {
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vgai-Key", XKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	return string(body), nil
}
