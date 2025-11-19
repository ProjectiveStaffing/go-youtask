package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kevino117/go-youtask/pkg/config"
	"github.com/kevino117/go-youtask/pkg/model"
	"github.com/kevino117/go-youtask/pkg/service"
)

// Handler es el entrypoint que Vercel invoca
func Handler(w http.ResponseWriter, r *http.Request) {
	// Lista de orígenes permitidos
	allowedOrigins := []string{
		"http://localhost:3000",
		"https://ai-assistant-one-liard.vercel.app",
		"https://ai-assistant-7kjd2dvzd-projective-staffings-projects.vercel.app",
		"https://ai-assistant-h2focl5c9-projective-staffings-projects.vercel.app",
	}

	// Obtener el origen de la petición
	origin := r.Header.Get("Origin")

	// Verificar si el origen está permitido
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}

	// Configurar otros headers CORS
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

	// Manejar preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Solo permitir POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed. Use POST",
		})
		return
	}

	// Parsear request
	var req model.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Validar que prompt no esté vacío
	if req.Prompt == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Prompt is required",
		})
		return
	}

	// Cargar configuración
	cfg := config.LoadConfig()

	// Llamar al servicio de Azure OpenAI
	response, err := service.GenerateTaskFromAzure(req.Prompt, cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
