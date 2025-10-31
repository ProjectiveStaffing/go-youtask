package model

type TaskRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type TaskResponse struct {
	Response ResponseData `json:"response"`
}

type ResponseData struct {
	TaskName       []string `json:"taskName"`
	PeopleInvolved []string `json:"peopleInvolved"`
	TaskCategory   []string `json:"taskCategory"`
	DateToPerform  string   `json:"dateToPerform"`
	ModelResponse  string   `json:"modelResponse"`
}
