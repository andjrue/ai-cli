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
		}
		providerModels[name] = providerCfg.Models
	}

	manager := models.Manager{
		providers,
		providerModels,
		cfg.DefaultProvider,
		cfg.DefaultModel,
	}
	app := ui.NewApp(&manager)
	if err := app.Run(); err != nil {
		log.Fatal("failed to start app")
	}

}
