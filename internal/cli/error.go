package cli

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ModelError struct {
	width, height int
	err           error
}

func NewModelError(err error) *ModelError {
	return &ModelError{err: err}
}

func (e *ModelError) Error() string {
	return e.err.Error()
}

func (e *ModelError) Init() tea.Cmd {
	return nil
}

func (e *ModelError) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.height = msg.Height
		e.width = msg.Width
	case tea.KeyMsg:
		return e, tea.Quit
	default:
		return e, nil
	}
	return e, nil
}

func (e *ModelError) View() string {
	return form(e.width, e.height, "Error occured! Press any key to quit", e.Error())
}
