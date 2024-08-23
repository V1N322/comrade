package comradelm

import (
	"fmt"

	lib "github.com/V1N322/comrade/lib"

)


type ComradeLM struct {
	Token               string
	Agent               string
	Context             []map[string]interface{}

	URL             string
	AutoContext     bool
}


func NewComradeLM(userURL string, token string, agent string, autoContext bool) *ComradeLM {
	return &ComradeLM{URL: userURL, Token: token, Agent: agent, AutoContext: autoContext}
}

func (comrade *ComradeLM) AddMessage(message string, role string) {
	comrade.Context = append(comrade.Context, map[string]interface{}{"role": role, "content": message})
}

func getStringContext(context []map[string]interface{}) string {
	var result string
	for _, message := range context {
		result += fmt.Sprintf("%s: %s\n", message["role"], message["content"])
	}
	return result
}

func (comrade *ComradeLM) ClearContext() {
	comrade.Context = []map[string]interface{}{}
}

func (comrade *ComradeLM) GetContext() []map[string]interface{} {
	return comrade.Context
}

func (comrade *ComradeLM) LoadContext(context []map[string]interface{}) {
	comrade.Context = context
}


func (comrade *ComradeLM) SendMessage(input string) (string, error) {
	if comrade.AutoContext {
		comrade.AddMessage(input, "user")
	}

	context := getStringContext(comrade.Context)

	req := lib.Request{
		ComradeAIToken: comrade.Token,
		Text:           context,
		AgentAddress:   comrade.Agent,
	}

	resp, err := lib.GetComradeAIResponse(req, comrade.URL)
	if err != nil {
		return "", fmt.Errorf("error getting response: %v", err)
	}

	if resp.Result != "success" {
		return "", fmt.Errorf("invalid response: %v", resp.Result)
	}

	content, ok := resp.Content.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid data structure")
	}

	lastTextOutput, ok := content["last_text_output"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid data structure")
	}

	text, ok := lastTextOutput["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid data structure")
	}

	if comrade.AutoContext {
		comrade.AddMessage(text, "assistant")
	}

	return text, nil
}
