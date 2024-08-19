package comradelm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bytes"

	"os"
)


type ComradeLM struct {
	Token               string
	Agent               string
	Context             []map[string]interface{}
}

type Response struct {
	ComradeAIToken string                   `json:"comradeAIToken"`
	Text           string                   `json:"text"`
	AgentAddress   string                   `json:"agentAddress"`
	RequestAgentConfig map[string]interface{} `json:"requestAgentConfig"`
}

func NewComradeLM(token string, agent string) *ComradeLM {
	return &ComradeLM{Token: token, Agent: agent}
}

func (comrade *ComradeLM) AddMessage(message string, role string) {
	comrade.Context = append(comrade.Context, map[string]interface{}{"role": role, "content": message})
}

func (comrade *ComradeLM) SendMessage(message string) (string, error) {
	AddMessage(comrade, message, "user")
}
