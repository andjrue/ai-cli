package components

import (
	"github.com/ai-cli/internal/models"
	"github.com/rivo/tview"
)

type ModelDropdown struct {
	MDropdown *tview.DropDown
	manager *models.Manager
	onProviderChange func(string, int)
	onModelChange func(string, int)
}

func NewModelDropdown(m *models.Manager, opc, omc func(string, int)) *ModelDropdown {
	return &ModelDropdown{
		MDropdown: tview.NewDropDown(),
		manager: m,
		onProviderChange: opc,
		onModelChange: omc,
	}
}

func (md *ModelDropdown) SetModelDropdown() {
	md.MDropdown = tview.NewDropDown()
	md.MDropdown.SetBorder(true).SetTitle("Model").SetTitleAlign(0)
}

func (md *ModelDropdown) SetModelOptions() *tview.DropDown {
	md.MDropdown.SetOptions(md.manager.GetModelsForCurrentProvider(), md.onModelChange).
		SetCurrentOption(0)
	
	return md.MDropdown
}