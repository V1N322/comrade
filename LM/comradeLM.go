package comradelm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"io"

	"bytes"

	httpClient "github.com/V1N322/httpUtils/httpConstant"
)


type ComradeLM struct {
	Token               string
	Agent               string
	Context             []map[string]interface{}

	URL             string
}

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


func NewComradeLM(userURL string, token string, agent string) *ComradeLM {
	return &ComradeLM{URL: userURL, Token: token, Agent: agent}
}

func (comrade *ComradeLM) AddMessage(message string, role string) {
	comrade.Context = append(comrade.Context, map[string]interface{}{"role": role, "content": message})
}

func (comrade *ComradeLM) getStringContext() string {
	var result string
	for _, message := range comrade.Context {
		result += fmt.Sprintf("%s: %s\n", message["role"], message["content"])
	}
	return result
}

func newPostRequest(URL string, jsonData []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s/get_agent_response/", URL)
	req, err := http.NewRequest(httpClient.POST, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	return req, nil
}


func (comrade *ComradeLM) SendMessage(message string) (string, error) {
	comrade.AddMessage(message, "user")

	context := comrade.getStringContext()

	request := Request{
		ComradeAIToken: comrade.Token,
		Text:           context,
		AgentAddress:   comrade.Agent,
		RequestAgentConfig: map[string]interface{}{},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "Something went wrong", err
	}

	req, err := newPostRequest(comrade.URL, jsonData)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := httpClient.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	var result ResponseData
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if result.Result == "success" {
		content, ok := result.Content.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid data structure")
		}
		lastTextOutput, ok := content["last_text_output"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid data structure")
		}
		responseText, ok := lastTextOutput["content"].(string)
		if !ok {
			return "", fmt.Errorf("invalid data structure")
		}

		comrade.AddMessage(responseText, "assistant")

		return responseText, nil
	} else {
		return fmt.Sprintf("%v", result.Content), nil
	}

	

	return "", nil
}
