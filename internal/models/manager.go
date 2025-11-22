package models

import (
	"fmt"
	"slices"
)

type Manager struct {
	Providers       map[string]Provider
	Models          map[string][]string
	CurrentProvider string
	CurrentModel    string
}

func NewManager(providers map[string]Provider, models map[string][]string, defaultProvider string, defaultModel string) *Manager {
	// Defaults will be loaded on each app restart
	return &Manager{
		Providers:       providers,
		Models:          models,
		CurrentProvider: defaultProvider,
		CurrentModel:    defaultModel,
	}
}

func (m *Manager) SwitchProvider(name string) error {
	if _, exists := m.Providers[name]; !exists {
		return fmt.Errorf("provider [%s] not found", name)
	}

	m.CurrentProvider = name

	if len(m.Models[name]) > 0 {
		m.CurrentModel = m.Models[name][0]
	}
	return nil
}

func (m *Manager) SwitchModel(model string) error {
	validModels := m.Models[m.CurrentProvider]

	found := slices.Contains(validModels, model)
	if found {
		m.CurrentModel = model
	}

	return fmt.Errorf("model [%s] not available from provider [%s]", model, m.CurrentProvider)
}

func (m *Manager) GetModelsForCurrentProvider() []string {
	return m.Models[m.CurrentProvider]
}

func (m *Manager) GetCurrentModel() string {
	return m.CurrentModel
}

func (m *Manager) GetCurrentProvider() Provider {
	return m.Providers[m.CurrentProvider]
}

func (m *Manager) GetProviderNames() []string {
	names := make([]string, 0, len(m.Providers))
	for name := range m.Providers {
		names = append(names, name)
	}

	return names
}
