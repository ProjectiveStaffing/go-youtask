package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kevino117/go-youtask/internal/config"
	"github.com/kevino117/go-youtask/internal/model"
	"github.com/kevino117/go-youtask/internal/service"
)

var router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://ai-assistant-6yceqcyen-projective-staffings-projects.vercel.app",
			"https://ai-assistant-one-liard.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))

	// Ruta principal para Vercel
	r.POST("/api/task", postTaskHandler)

	router = r
}

func postTaskHandler(c *gin.Context) {
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

// Handler es el entrypoint que Vercel invoca
func Handler(w http.ResponseWriter, r *http.Request) {
	// Vercel espera que manejemos cualquier ruta que llegue
	// Reescribimos la ruta para que coincida con nuestro router
	if r.URL.Path == "/api/task" || r.URL.Path == "/api/task/" {
		router.ServeHTTP(w, r)
		return
	}

	// Si no es la ruta correcta, devolvemos 404
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Not found. Use POST /api/task",
	})
}
