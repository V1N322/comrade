package comradeembedding

import (
	"fmt"

	lib "comrade/lib"
)


type ComradeEmbedding struct {
	Token               string
	Agent               string

	URL             	string
}


func NewComradeEmbedding(userURL string, token string, agent string) *ComradeEmbedding {
	return &ComradeEmbedding{URL: userURL, Token: token, Agent: agent}
}


func (comrade *ComradeEmbedding) EmbedText(message string) (string, error) {

	fmt.Printf("Sending message '%s' to Comrade Embedding\n", message)

	request := lib.Request{
		ComradeAIToken: comrade.Token,
		Text:           message,
		AgentAddress:   comrade.Agent,
		RequestAgentConfig: map[string]interface{}{},
	}

	fmt.Println("Sending request to Comrade API")
	result, err := lib.GetComradeAIResponse(request, comrade.URL)
	if err != nil {
		fmt.Printf("Error getting response: %v\n", err)
		return "", fmt.Errorf("error getting response: %v", err)
	}

	fmt.Println("Processing response")
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


		// fmt.Printf("Returning response: %s\n", lastTextOutput)
		return responseText, nil
	}

	fmt.Println("Invalid response")
	return "", fmt.Errorf("invalid response")
}