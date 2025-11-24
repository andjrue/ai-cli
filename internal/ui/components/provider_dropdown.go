// Package components creates components for the app to use
package components

import (
	"github.com/ai-cli/internal/models"
	"github.com/rivo/tview"
)

type ProviderDropdown struct {
	PDropdown         *tview.DropDown
	manager           *models.Manager
	onProviderChanged func(string, int)
	onModelChanged    func(string, int)
}

func NewProviderDropdown(manager *models.Manager, onProviderChanged, onModelChanged func(string, int)) *ProviderDropdown {
	return &ProviderDropdown{
		PDropdown:         tview.NewDropDown(),
		manager:           manager,
		onProviderChanged: onProviderChanged,
		onModelChanged:    onModelChanged,
	}

}

func (pd *ProviderDropdown) SetProviderDropdown() {
	pd.PDropdown = tview.NewDropDown()
	pd.PDropdown.SetBorder(true).SetTitle("Provider").SetTitleAlign(0)
}

func (pd *ProviderDropdown) SetProviderOptions() *tview.DropDown {
	pd.PDropdown.SetOptions(pd.manager.GetProviderNames(), nil)
	
	currentProvider := pd.manager.CurrentProvider
	for i, name := range pd.manager.GetProviderNames() {
		if name == currentProvider {
			pd.PDropdown.SetCurrentOption(i)
			break
		}
	}
	
	pd.PDropdown.SetSelectedFunc(pd.onProviderChanged)

	return pd.PDropdown
}
