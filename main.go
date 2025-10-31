package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/kevino117/go-youtask/internal/config"
	"github.com/kevino117/go-youtask/internal/handler"
)

func main() {
	cfg := config.LoadConfig()
	r := gin.Default()

	config := cors.Config{
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

	r.Use(cors.New(config))

	r.POST("/youtask/api/v0/task", handler.PostTaskHandler)

	log.Printf("🚀 Servidor escuchando en http://localhost:%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
