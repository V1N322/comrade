package comradelm

import (
	"fmt"

	lib "comrade/lib"

)


type ComradeLM struct {
	Token               string
	Agent               string
	Context             []map[string]interface{}

	URL             string
	AutoContext     bool
}


func NewComradeLM(userURL string, token string, agent string) *ComradeLM {
	return &ComradeLM{URL: userURL, Token: token, Agent: agent}
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



func (comrade *ComradeLM) SendMessage(message string) (string, error) {
	
	if comrade.AutoContext {
		comrade.AddMessage(message, "user")
	}

	context := getStringContext(comrade.Context)

	request := lib.Request{
		ComradeAIToken: comrade.Token,
		Text:           context,
		AgentAddress:   comrade.Agent,
		RequestAgentConfig: map[string]interface{}{},
	}


	result, err := lib.GetComradeAIResponse(request, comrade.URL)
	if err != nil {
		return "", fmt.Errorf("error getting response: %v", err)
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

		if comrade.AutoContext {
			comrade.AddMessage(responseText, "assistant")
		}

		return responseText, nil

	} else {
		return fmt.Sprintf("%v", result.Content), nil
	}

	return "", nil
}
