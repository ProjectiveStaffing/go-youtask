package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kevino117/go-youtask/internal/config"
	"github.com/kevino117/go-youtask/internal/model"
	"github.com/kevino117/go-youtask/internal/service"
)

// Handler es el entrypoint que Vercel invoca
func Handler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
