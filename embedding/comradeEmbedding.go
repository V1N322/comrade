package comradeembedding

import (
	"fmt"

	lib "comrade/lib"

	"encoding/json"
)


type ComradeEmbedding struct {
	Token               string
	Agent               string

	URL             	string
}

type EmbeddingResult struct {
	Model      string    `json:"model"`
	Embeddings []float64 `json:"embeddings"`
}

func NewComradeEmbedding(userURL string, token string, agent string) *ComradeEmbedding {
	return &ComradeEmbedding{URL: userURL, Token: token, Agent: agent}
}

type ComradeAPIResponse struct {
	Result  string                 `json:"result"`
	Content map[string]interface{} `json:"content"`
}

func (comrade *ComradeEmbedding) EmbedText(message string) ([]float64, error) {
	fmt.Printf("Sending message '%s' to Comrade Embedding\n", message)

	request := lib.Request{
		ComradeAIToken: comrade.Token,
		Text:           message,
		AgentAddress:   comrade.Agent,
		RequestAgentConfig: map[string]interface{}{},
	}

	result, err := lib.GetComradeAIResponse(request, comrade.URL)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %v", err)
	}

	if result.Result == "success" {
		content, ok := result.Content.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data structure")
		}

		contentList, ok := content["last_text_output"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data structure")
		}
		
		modelsList, ok := contentList["content"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid data structure")
		}

		var embeddingsResult []EmbeddingResult

		err := json.Unmarshal([]byte(modelsList), &embeddingsResult)
		if err != nil {
			return nil, fmt.Errorf("error parsing JSON: %v", err)
		}

		for _, model := range embeddingsResult {

			if model.Model == "all-mpnet-base-v2" {

				var floatEmbeddings []float64
				floatEmbeddings = append(floatEmbeddings, model.Embeddings...)

				return floatEmbeddings, nil
			}
		}

		return nil, fmt.Errorf("model 'all-mpnet-base-v2' not found")
	}

	return nil, fmt.Errorf("error getting response: %s", result.Result)
}