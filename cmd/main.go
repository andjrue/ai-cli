package main

import (
	"log"

	"github.com/ai-cli/internal/config"
	"github.com/ai-cli/internal/models"
	"github.com/ai-cli/internal/ui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: %w", err)
	}
	
	var provider models.Provider
	if cfg.Models.OpenAI.APIKey != "" {
		provider = models.NewOpenAIProvider(cfg.Models.OpenAI.APIKey)
	} else {
		log.Fatal("no api key configured for open ai")
	}
	
	app := ui.NewApp(provider)
	if err := app.Run(); err != nil {
		log.Fatal("failed to start app")
	}
	
}