package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AzureConfig struct {
	Endpoint       string
	ApiKey         string
	DeploymentName string
	ApiVersion     string
	Port           string
}

func LoadConfig() AzureConfig {
	_ = godotenv.Load()

	cfg := AzureConfig{
		Endpoint:       os.Getenv("AZURE_OPENAI_ENDPOINT"),
		ApiKey:         os.Getenv("AZURE_OPENAI_API_KEY"),
		DeploymentName: os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME"),
		ApiVersion:     os.Getenv("AZURE_OPENAI_API_VERSION"),
		Port:           os.Getenv("PORT"),
	}

	if cfg.Endpoint == "" || cfg.ApiKey == "" || cfg.DeploymentName == "" || cfg.ApiVersion == "" {
		log.Fatal("❌ Missing enviroment variables for Azure OpenAI.")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}
