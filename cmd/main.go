package main

import (
	"log"

	"github.com/ai-cli/internal/config"
	"github.com/ai-cli/internal/logger"
	"github.com/ai-cli/internal/models"
	"github.com/ai-cli/internal/ui"
)

func main() {

	if err := logger.Init("debug.log"); err != nil {
		panic(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: %w", err)
	}

	providers := make(map[string]models.Provider)
	providerModels := make(map[string][]string)

	for name, providerCfg := range cfg.Models {
		switch name {
		case "openai":
			providers[name] = models.NewOpenAIProvider(providerCfg.APIKey)
			
		case "anthropic":
			providers[name] = models.NewAnthropicProvider(providerCfg.APIKey)
		}
		providerModels[name] = providerCfg.Models
	}

	manager := models.NewManager(providers, providerModels, cfg.DefaultProvider, cfg.DefaultModel)
	logger.Log.Printf("------------------------")
	logger.Log.Printf("Default provider from config: %s", cfg.DefaultProvider)
	logger.Log.Printf("Default model from config: %s", cfg.DefaultModel)
	logger.Log.Printf("Manager current provider: %s", manager.CurrentProvider)
	logger.Log.Printf("Manager current model: %s", manager.CurrentModel)
	logger.Log.Printf("Available providers: %s", manager.GetProviderNames())
	logger.Log.Printf("Models for [%s]: %v", manager.CurrentProvider, manager.GetModelsForCurrentProvider())
	logger.Log.Printf("------------------------")
	app := ui.NewApp(manager)
	if err := app.Run(); err != nil {
		log.Fatal("failed to start app")
	}

}
