package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sebasttiano/Owl/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrInitMainBoard = errors.New("failed to init main board")
)

type resType int

const margin = 4

func (c *CLI) StartMainBoard(ctx context.Context) error {

	mainBoard := NewMainBoard(c)
	if err := mainBoard.initLists(ctx); err != nil {
		logger.Log.Error("failed to init board columns", zap.Error(err))
		return ErrInitMainBoard
	}

	_, err := tea.NewProgram(
		mainBoard,
		tea.WithAltScreen(),
		tea.WithContext(ctx)).Run()

	if err != nil {
		return ErrInitMainBoard
	}
	return nil
}

const (
	credType resType = iota
	cardType
	textType
)

type MainBoard struct {
	help      help.Model
	cols      []column
	loaded    bool
	focused   resType
	list      list.Model
	cancelled bool
	cli       *CLI
}

func (m *MainBoard) initLists(ctx context.Context) error {

	// Init cred, card, text types
	m.cols = []column{
		newColumn(credType),
		newColumn(cardType),
		newColumn(textType),
	}

	m.cols[credType].list.Title = "Credentials"
	m.cols[cardType].list.Title = "Bank cards"
	m.cols[textType].list.Title = "Notes"

	resp, err := m.cli.Client.Text.GetAllTexts(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	for i, text := range resp.GetTexts() {
		item := ResourceItem{resType: textType, title: fmt.Sprintf("ID: %d", text.Id), description: text.Description}
		m.cols[textType].list.InsertItem(i, item)

	}
	return nil
}

func (m *MainBoard) UpdateColumns() {

}
func NewMainBoard(cli *CLI) *MainBoard {
	help := help.New()
	help.ShowAll = true
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return &MainBoard{help: help, list: defaultList, cli: cli}
}

func (m *MainBoard) Init() tea.Cmd {
	return nil
}

func (m *MainBoard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		for i := 0; i < len(m.cols); i++ {
			var res tea.Model
			res, cmd = m.cols[i].Update(msg)
			m.cols[i] = res.(column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.cancelled = true
			return m, tea.Quit
			//case key.Matches(msg, keys.Enter):
			//	switch option := s.list.SelectedItem().(type) {
			//	case signInModel:
			//		return option.Update(nil)
			//	case signUpModel:
			//		return option.Update(nil)
			//	}
		}
	}
	return m, nil

}

func (m *MainBoard) View() string {
	if m.cancelled {
		return ""
	}
	if !m.loaded {
		return "LOADING!"
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[credType].View(),
		m.cols[cardType].View(),
		m.cols[textType].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keys))
}
