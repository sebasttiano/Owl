package cli

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const margin = 4

const (
	signIn state = iota
	signUp
	credType
	cardType
	textType
)

type MainModel struct {
	width, height int
	help          help.Model
	loaded        bool
	focused       state
	list          list.Model
	quitting      bool
}

func NewMainModel() *MainModel {
	help := help.New()
	help.ShowAll = true
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return &MainModel{help: help, list: defaultList}
}

func (m *MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
	//switch msg := msg.(type) {
	//case tea.WindowSizeMsg:
	//	var cmd tea.Cmd
	//	var cmds []tea.Cmd
	//	m.height = msg.Height
	//	m.width = msg.Width
	//	m.help.Width = msg.Width - margin
	//}

}

func (m *MainModel) View() string {
	return ""
	//if m.quitting {
	//	return ""
	//}
	//if !m.loaded {
	//	return "LOADING!"
	//}
	//output := lipgloss.JoinHorizontal(
	//	lipgloss.Left,
	//	m.View(),
	//	m.help.View(keys))

	//return output
}
