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

func (comrade *ComradeEmbedding) EmbedText(input string) ([]float64, error) {
	req := lib.Request{
		ComradeAIToken: comrade.Token,
		Text:           input,
		AgentAddress:   "Embeddings",
		RequestAgentConfig: map[string]interface{}{},
	}

	resp, err := lib.GetComradeAIResponse(req, comrade.URL)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}

	if resp.Result != "success" {
		return nil, fmt.Errorf("error getting response: %s", resp.Result)
	}

	var models []EmbeddingResult
	modelsJSON, ok := resp.Content.(map[string]interface{})["last_text_output"].(map[string]interface{})["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid data structure")
	}

	if err := json.Unmarshal([]byte(modelsJSON), &models); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	for _, model := range models {
		if model.Model == comrade.Agent {
			return model.Embeddings, nil
		}
	}

	return nil, fmt.Errorf("model %s not found", comrade.Agent)
}
