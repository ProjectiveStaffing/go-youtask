package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kevino117/go-youtask/internal/config"
	"github.com/kevino117/go-youtask/internal/model"
)

const systemPrompt = `You are an AI agent specialized in task management. 
Your job is to extract from the user's message:
1. The name(s) of the task(s).
2. The people involved.
3. The category of the task ("Family", "Work", or "Other").
4. The date or time to perform the task.
5. A short natural language summary called "modelResponse".

You must answer **only** with a valid JSON object that matches this exact structure.
Do not include explanations, markdown, or any extra characters outside of the JSON.

Expected JSON structure (use JSON array brackets for lists):

{
  "taskName": ["string"],
  "peopleInvolved": ["string"],
  "taskCategory": ["Family" | "Work" | "Other"],
  "dateToPerform": "string",
  "modelResponse": "string"
}`

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages         []ChatMessage `json:"messages"`
	MaxTokens        int           `json:"max_tokens"`
	Temperature      float32       `json:"temperature"`
	TopP             float32       `json:"top_p"`
	FrequencyPenalty float32       `json:"frequency_penalty"`
	PresencePenalty  float32       `json:"presence_penalty"`
}

type AzureResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateTaskFromAzure(prompt string, cfg config.AzureConfig) (model.TaskResponse, error) {
	reqBody := ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		MaxTokens:        1024,
		Temperature:      0.7,
		TopP:             0.95,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	body, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		cfg.Endpoint, cfg.DeploymentName, cfg.ApiVersion)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return model.TaskResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", cfg.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model.TaskResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return model.TaskResponse{}, fmt.Errorf("error: %s", string(b))
	}

	var azureResp AzureResponse
	if err := json.NewDecoder(resp.Body).Decode(&azureResp); err != nil {
		return model.TaskResponse{}, err
	}

	content := azureResp.Choices[0].Message.Content

	log.Printf("🔍 Azure OpenAI raw response: %s\n", content)

	var parsed model.ResponseData
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return model.TaskResponse{}, fmt.Errorf("error parsing model response: %v", err)
	}

	return model.TaskResponse{Response: parsed}, nil
}
