package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kevino117/go-youtask/pkg/config"
	"github.com/kevino117/go-youtask/pkg/model"
)

const systemPrompt = `You are an AI agent specialized in personal task management.

Your job is to extract structured information from the user's message. 
The user might describe something they need to do (task), a recurring activity (habit), or a bigger goal (project).

You must analyze the message and determine:

1. "taskName": the main name(s) of the task(s) mentioned.
2. "peopleInvolved": all people explicitly mentioned in the message.
3. "taskCategory": one or more of ["Family", "Work", "Other"] depending on context.
4. "dateToPerform": the specific day, date, or time extracted from the message (if any).
5. "modelResponse": a natural language short summary confirming the action.
6. "itemType": classify the message as one of ["Task", "Habit", "Project"].
7. "assignedTo": identify who is responsible for doing it.
   - If the message implies the user will do it (e.g., "I have to..."), set it to "User".
   - If the message implies someone else will do it (e.g., "My mom has to..."), set it to that person’s name.

You must output ONLY a valid JSON object matching this exact structure, with no markdown or explanations.

Expected JSON structure:

{
  "taskName": ["string"],
  "peopleInvolved": ["string"],
  "taskCategory": ["Family" | "Work" | "Other"],
  "dateToPerform": "string",
  "modelResponse": "string",
  "itemType": ["Task" | "Habit" | "Project"],
  "assignedTo": "string"
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

	var parsed model.ResponseData
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return model.TaskResponse{}, fmt.Errorf("error parsing model response: %v", err)
	}

	return model.TaskResponse{Response: parsed}, nil
}
