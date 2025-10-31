package service

import "github.com/kevino117/go-youtask/internal/model"

func GenerateTaskResponse(prompt string) model.TaskResponse {
	return model.TaskResponse{
		Response: model.ResponseData{
			TaskName:       []string{"Clean the mirror in the hallway"},
			PeopleInvolved: []string{"My brother", "My cousin"},
			TaskCategory:   []string{"Family"},
			DateToPerform:  "Tomorrow morning",
			ModelResponse:  "You have been assigned to clean the mirror in the hallway tomorrow morning by your brother and cousin. Task category: Family.",
		},
	}
}
