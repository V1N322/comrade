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

func newPostRequest(url string, jsonData []byte) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/get_agent_response/", url)
	request, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	return request, nil
}


func GetComradeAIResponse(request Request, url string) (ResponseData, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error creating JSON: %w", err)
	}

	req, err := newPostRequest(url, data)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := httpClient.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResponseData{}, fmt.Errorf("invalid response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error reading response body: %w", err)
	}

	var result ResponseData
	if err := json.Unmarshal(body, &result); err != nil {
		return ResponseData{}, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}
