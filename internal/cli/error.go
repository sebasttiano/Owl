package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type ErrorModel struct {
	width, height int
	err           error
}

func NewErrorModel(err error) ErrorModel {
	return ErrorModel{err: err}
}

func (e ErrorModel) Error() string {
	return e.err.Error()
}

func (e ErrorModel) Init() tea.Cmd {
	return nil
}

func (e ErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.height = msg.Height
		e.width = msg.Width
	default:
		time.Sleep(3 * time.Second)
		return e, tea.Quit
	}
	return e, nil
}

func (e ErrorModel) View() string {
	return form(e.width, e.height, "Error occured!", e.Error())
}
