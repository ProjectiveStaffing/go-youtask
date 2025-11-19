# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**go-youtask** is a Go-based API service that uses Azure OpenAI to extract structured task information from natural language prompts. The service analyzes user messages to identify tasks, habits, or projects and returns structured JSON with task details, involved people, categories, dates, and assignments.

## Development Commands

### Running the Application
```bash
go run main.go
```
Server starts on port 8080 by default (or PORT env var).

### Building
```bash
go build -o youtask
```

### Installing Dependencies
```bash
go mod download
go mod tidy
```

### Testing
```bash
go test ./...
```

## Deployment to Vercel

### Prerequisites
1. GitHub repository connected to Vercel
2. Environment variables configured in Vercel dashboard

### Vercel Environment Variables Setup
In your Vercel project settings, add these environment variables:
- `AZURE_OPENAI_ENDPOINT`
- `AZURE_OPENAI_API_KEY`
- `AZURE_OPENAI_DEPLOYMENT_NAME`
- `AZURE_OPENAI_API_VERSION`

### Deployment Process
1. Push changes to your GitHub repository
2. Vercel automatically detects and deploys the Go serverless function
3. The API will be available at: `https://your-project.vercel.app/api/task`

### API Endpoint in Production
**POST https://your-project.vercel.app/api/task**
- Same request/response format as local development
- CORS pre-configured for allowed origins

### Local vs Vercel
- **Local**: `main.go` runs a persistent Gin server
- **Vercel**: `api/task.go` runs as a serverless function (cold starts possible)
- Both use the same handler logic and configuration

## Environment Configuration

Required environment variables (loaded from `.env`):
- `AZURE_OPENAI_ENDPOINT` - Azure OpenAI endpoint URL
- `AZURE_OPENAI_API_KEY` - Azure OpenAI API key
- `AZURE_OPENAI_DEPLOYMENT_NAME` - Deployment name for the model
- `AZURE_OPENAI_API_VERSION` - API version (e.g., "2024-02-15-preview")
- `PORT` - Server port (defaults to 8080 if not set)

Missing any of these (except PORT) will cause the application to exit on startup.

## Architecture

### Layer Structure

The codebase follows a clean layered architecture:

**Handler Layer** ([internal/handler/](internal/handler/))
- HTTP request/response handling via Gin
- Request validation and JSON binding
- Error responses with appropriate status codes

**Service Layer** ([internal/service/](internal/service/))
- Business logic and external service integration
- `GenerateTaskFromAzure` calls Azure OpenAI API with system prompt
- Marshals requests, unmarshals structured JSON responses from the model

**Model Layer** ([internal/model/](internal/model/))
- Data structures for requests/responses
- `TaskRequest`: incoming user prompt
- `TaskResponse`: structured response with `ResponseData`
- `ResponseData`: taskName, peopleInvolved, taskCategory, dateToPerform, modelResponse, itemType, assignedTo

**Config Layer** ([internal/config/](internal/config/))
- Environment variable loading via godotenv
- Configuration validation at startup

### API Endpoints

**POST /youtask/api/v0/task**
- Accepts JSON: `{"prompt": "user's natural language message"}`
- Returns structured task data extracted by Azure OpenAI
- Handler: `PostTaskHandler` in [internal/handler/task_handler.go](internal/handler/task_handler.go)

### Deployment Context

The codebase contains two entry points:

1. **main.go** - Standard Go server with Gin router (local development)
2. **api/task.go** - Vercel serverless function with `Handler(w, r)` export

The serverless function uses `init()` to set up the Gin router once per cold start, then `Handler()` serves each request through that router.

### Azure OpenAI Integration

The system prompt in [internal/service/openai_sevice.go](internal/service/openai_sevice.go) defines strict JSON output formatting. The model must classify messages into:
- **itemType**: Task, Habit, or Project
- **taskCategory**: Family, Work, or Other
- **assignedTo**: "User" if user will do it, otherwise the person's name

The service expects valid JSON directly from the model (no markdown wrappers).

## Code Patterns

- Configuration is loaded on demand in handlers (`config.LoadConfig()`), not cached globally
- HTTP client is created per-request in `GenerateTaskFromAzure` (not reused)
- Error handling returns early with appropriate HTTP status codes
- CORS is configured for specific production Vercel URLs and localhost:3000
