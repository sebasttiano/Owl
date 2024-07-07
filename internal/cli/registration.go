package cli

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type registrationModel struct {
	cancelled     bool
	width, height int

	username textinput.Model
	password textinput.Model
}

func (r registrationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (r registrationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		usernameCmd tea.Cmd
		passwordCmd tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.height = msg.Height
		r.width = msg.Width
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			switch {
			case r.username.Focused():
				if len(r.username.Value()) < 1 {
					return r, textinput.Blink
				}
				r.username.Blur()
				r.password.Focus()
			case r.password.Focused():
				if len(r.password.Value()) < 1 {
					return r, textinput.Blink
				}
				r.password.Blur()
				return r, tea.Quit
			}
			if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
				r.cancelled = true
				return r, tea.Quit
			}
		}
	}
	return r, tea.Batch(usernameCmd, passwordCmd)
}

// View implements tea.Model.
func (r registrationModel) View() string {
	return form(
		r.width, r.height,
		"Register to Owl",
		lipgloss.JoinVertical(
			lipgloss.Left,
			r.username.View(),
			strings.Repeat(" ", 64),
			r.password.View(),
		),
	)
}

func newRegistrationModel() registrationModel {
	m := registrationModel{
		cancelled: false,
		username:  textinput.New(),
		password:  textinput.New(),
	}
	m.username.CharLimit = 32
	m.username.Prompt = "Username: "
	m.username.Placeholder = "type your username..."

	m.password.CharLimit = 32
	m.password.Prompt = "Password: "
	m.password.EchoMode = textinput.EchoPassword
	m.password.Placeholder = "type your password..."

	m.username.Focus()
	return m
}
