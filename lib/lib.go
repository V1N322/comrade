package lib

import (
	"bytes"
	"fmt"
	"net/http"

	"encoding/json"

	"io"

	httpClient "github.com/V1N322/httpUtils/httpConstant"
)

type Request struct {
	ComradeAIToken string                   `json:"comradeAIToken"`
	Text           string                   `json:"text"`
	AgentAddress   string                   `json:"agentAddress"`
	RequestAgentConfig map[string]interface{} `json:"requestAgentConfig"`
}

type ResponseData struct {
	Result  string      `json:"result"`
	Content interface{} `json:"content"`
}

func NewPostRequest(URL string, jsonData []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s/get_agent_response/", URL)
	req, err := http.NewRequest(httpClient.POST, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	return req, nil
}

func GetComradeAIResponse(request Request, URL string) (ResponseData, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error creating json: %v", err)
	}

	req, err := NewPostRequest(URL, jsonData)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := httpClient.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResponseData{}, fmt.Errorf("invalid response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error reading response body: %v", err)
	}
	var result ResponseData
	if err := json.Unmarshal(body, &result); err != nil {
		return ResponseData{}, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}
