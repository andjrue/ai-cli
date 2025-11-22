// Package ui provides all required UI items and styling
package ui

import (
	"context"
	"fmt"

	"github.com/ai-cli/internal/models"
	"github.com/ai-cli/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	tviewApp         *tview.Application
	manager          *models.Manager
	providerSelector *components.ProviderDropdown
	modelSelector    *components.ModelDropdown
	output           *tview.TextView
	input            *tview.InputField
	layout           *tview.Flex
}

func NewApp(manager *models.Manager) *App {
	app := &App{
		tviewApp: tview.NewApplication(),
		manager:  manager,
	}

	app.setupUI()
	return app
}

func (a *App) setupUI() {

	// Init provider selector
	a.providerSelector = components.NewProviderDropdown(a.manager, a.onProviderChanged, a.onModelChanged)
	a.providerSelector.SetProviderDropdown()

	// Init model selector
	a.modelSelector = components.NewModelDropdown(a.manager, a.onProviderChanged, a.onModelChanged)
	a.modelSelector.SetModelDropdown()

	providerOptions := a.providerSelector.SetProviderOptions()
	modelOptions := a.modelSelector.SetModelOptions()

	a.output = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			a.output.ScrollToEnd()
			a.tviewApp.Draw()
		})
	a.output.SetBorder(true).SetTitle("Response")

	a.input = tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).
		SetDoneFunc(a.handleSubmit).
		SetFieldStyle(tcell.StyleDefault.
			Background(tcell.ColorDefault).
			Foreground(tcell.ColorDefault))
	a.input.SetBorder(true).SetTitle("Prompt")

	selectorRow := tview.NewFlex().
		AddItem(providerOptions, 0, 1, false).
		AddItem(modelOptions, 0, 1, false)

	a.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(selectorRow, 3, 0, false).
		AddItem(a.output, 0, 6, false).
		AddItem(a.input, 0, 2, true)

	a.tviewApp.SetRoot(a.layout, true)

	// Styling - Shout out claude for handling the important parts of this project
	tview.Styles.PrimitiveBackgroundColor = tcell.NewRGBColor(17, 17, 27)    // Darker base
	tview.Styles.ContrastBackgroundColor = tcell.NewRGBColor(30, 30, 46)     // Darker secondary
	tview.Styles.MoreContrastBackgroundColor = tcell.NewRGBColor(40, 40, 60) // Dropdown background
	tview.Styles.PrimaryTextColor = tcell.NewRGBColor(205, 214, 244)         // Soft white text
	tview.Styles.BorderColor = tcell.NewRGBColor(88, 91, 112)                // Muted grey borders
	tview.Styles.TitleColor = tcell.NewRGBColor(180, 190, 254)               // Soft blue for titles

	// Provider Dropdown
	a.providerSelector.PDropdown.SetBorderColor(tcell.NewRGBColor(88, 91, 112))
	a.providerSelector.PDropdown.SetTitleColor(tcell.NewRGBColor(235, 188, 186))
	a.providerSelector.PDropdown.SetBackgroundColor(tcell.NewRGBColor(17, 17, 27))
	a.providerSelector.PDropdown.SetFieldBackgroundColor(tcell.NewRGBColor(30, 30, 46))
	a.providerSelector.PDropdown.SetFieldTextColor(tcell.NewRGBColor(205, 214, 244))
	a.providerSelector.PDropdown.SetFocusFunc(func() {
		a.providerSelector.PDropdown.SetBorderColor(tcell.NewRGBColor(235, 160, 172)) // Rose pink
	})
	a.providerSelector.PDropdown.SetBlurFunc(func() {
		a.providerSelector.PDropdown.SetBorderColor(tcell.NewRGBColor(88, 91, 112))
	})

	// Model Selector
	a.modelSelector.MDropdown.SetBorderColor(tcell.NewRGBColor(88, 91, 112))
	a.modelSelector.MDropdown.SetTitleColor(tcell.NewRGBColor(235, 188, 186))
	a.modelSelector.MDropdown.SetBackgroundColor(tcell.NewRGBColor(17, 17, 27))
	a.modelSelector.MDropdown.SetFieldBackgroundColor(tcell.NewRGBColor(30, 30, 46))
	a.modelSelector.MDropdown.SetFieldTextColor(tcell.NewRGBColor(205, 214, 244))
	a.modelSelector.MDropdown.SetFocusFunc(func() {
		a.modelSelector.MDropdown.SetBorderColor(tcell.NewRGBColor(235, 160, 172))
	})
	a.modelSelector.MDropdown.SetBlurFunc(func() {
		a.modelSelector.MDropdown.SetBorderColor(tcell.NewRGBColor(88, 91, 112))
	})

	// Output TextView
	a.output.SetBorderColor(tcell.NewRGBColor(88, 91, 112))
	a.output.SetTitleColor(tcell.NewRGBColor(180, 190, 254)) // Soft blue
	a.output.SetTextColor(tcell.NewRGBColor(205, 214, 244))  // Main text
	a.output.SetFocusFunc(func() {
		a.output.SetBorderColor(tcell.NewRGBColor(166, 227, 161)) // Green accent
	})
	a.output.SetBlurFunc(func() {
		a.output.SetBorderColor(tcell.NewRGBColor(88, 91, 112))
	})

	// Input Field
	a.input.SetBorderColor(tcell.NewRGBColor(88, 91, 112))
	a.input.SetTitleColor(tcell.NewRGBColor(249, 226, 175)) // Warm yellow
	a.input.SetLabelColor(tcell.NewRGBColor(249, 226, 175))
	a.input.SetFieldTextColor(tcell.NewRGBColor(205, 214, 244))
	a.input.SetFieldBackgroundColor(tcell.NewRGBColor(30, 30, 46))
	a.input.SetFocusFunc(func() {
		a.input.SetBorderColor(tcell.NewRGBColor(249, 226, 175)) // Yellow accent
	})
	a.input.SetBlurFunc(func() {
		a.input.SetBorderColor(tcell.NewRGBColor(88, 91, 112))
	})
}

func (a *App) handleSubmit(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	prompt := a.input.GetText()
	if prompt == "" {
		return
	}

	a.input.SetText("")

	fmt.Fprintf(a.output, "\n[#f9e2af]Prompt:[#cdd6f4] %s\n\n", prompt)
	fmt.Fprintf(a.output, "[#a6e3a1]Assistant:[#cdd6f4]\n")

	req := models.Request{
		Model: a.manager.GetCurrentModel(),
		Messages: []models.Message{
			{Role: "user", Content: prompt},
		},
	}

	go a.streamResponse(req)
}

func (a *App) streamResponse(req models.Request) {
	ctx := context.Background()
	provider := a.manager.GetCurrentProvider()
	respChan, err := provider.Stream(ctx, req)
	if err != nil {
		a.tviewApp.QueueUpdateDraw(func() {
			fmt.Fprintf(a.output, "\n[#f38ba8]Error: %v[#cdd6f4]\n", err)
		})
	}

	for resp := range respChan {
		switch resp.Type {
		case models.ResponseTypeText:
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.output, resp.Content)
			})

		case models.ResponseTypeError:
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.output, "\n[red]Error streaming response: %v[white]", resp.Error)
			})

		case models.ResponseTypeDone:
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.output, "\n[#a6e3a1]Response complete[#cdd6f4]")
			})
		}
	}
}

func (a *App) getProviderNames() []string {
	names := make([]string, 0)
	for name := range a.manager.Providers {
		names = append(names, name)
	}

	return names
}

func (a *App) onProviderChanged(text string, index int) {
	a.manager.SwitchProvider(text)
	models := a.manager.GetModelsForCurrentProvider()
	a.modelSelector.MDropdown.SetOptions(models, a.onModelChanged)
	a.modelSelector.MDropdown.SetCurrentOption(0)
}

func (a *App) onModelChanged(text string, index int) {
	a.manager.SwitchModel(text)
}

func (a *App) Run() error {
	return a.tviewApp.Run()
}
