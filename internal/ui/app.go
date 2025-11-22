// Package ui provides all required UI items and styling
package ui

import (
	"context"
	"fmt"

	"github.com/ai-cli/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	tviewApp *tview.Application
	provider models.Provider
	output *tview.TextView
	input *tview.InputField
	layout *tview.Flex
}

func NewApp(provider models.Provider) *App {
	app := &App{
		tviewApp: tview.NewApplication(),
		provider: provider,
	}
	
	app.setupUI()
	return app
}

func (a *App) setupUI() {
	a.output = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
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
	
	a.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.output, 0, 8, false).
		AddItem(a.input, 0, 2, true)
	
	a.tviewApp.SetRoot(a.layout, true)
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
	
	fmt.Fprintf(a.output, "\n[yellow]Prompt:[white] %s\n\n", prompt)
	fmt.Fprintf(a.output, "[cyan]Assistant:[white]\n")
	
	req := models.Request{
		Model: "gpt-4o",
		Messages: []models.Message{
			{Role: "user", Content: prompt},
		},
	}
	
	go a.streamResponse(req)
}

func (a *App) streamResponse(req models.Request) {
	ctx := context.Background()
	respChan, err := a.provider.Stream(ctx, req)
	if err != nil {
		a.tviewApp.QueueUpdateDraw(func() {
			fmt.Fprintf(a.output, "\n[red]Error: %v[white]\n", err)
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
					fmt.Fprintf(a.output, "\n[green]Response complete[white]")
				})
		}
	}
}

func (a *App) Run() error {
	return a.tviewApp.Run()
}