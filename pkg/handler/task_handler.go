package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevino117/go-youtask/pkg/config"
	"github.com/kevino117/go-youtask/pkg/model"
	"github.com/kevino117/go-youtask/pkg/service"
)

func PostTaskHandler(c *gin.Context) {
	var req model.TaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cfg := config.LoadConfig()

	response, err := service.GenerateTaskFromAzure(req.Prompt, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
