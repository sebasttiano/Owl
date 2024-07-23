package cli

import (
	"context"
	"errors"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrInitMainBoard = errors.New("failed to init main board")
)

var mainBoard *MainBoard

type resType int

func (s resType) getNext() resType {
	if s == textType {
		return credType
	}
	return s + 1
}

func (s resType) getPrev() resType {
	if s == credType {
		return textType
	}
	return s - 1
}

const margin = 4

func (c *CLI) StartMainBoard(ctx context.Context) error {

	mainBoard = NewMainBoard(ctx, c)
	if err := mainBoard.initLists(); err != nil {
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
	ctx       context.Context
	help      help.Model
	cols      []column
	loaded    bool
	focused   resType
	list      list.Model
	cancelled bool
	cli       *CLI
}

func (m *MainBoard) initLists() error {

	// Init cred, card, text types
	m.cols = []column{
		newColumn(m.ctx, credType, m.cli),
		newColumn(m.ctx, cardType, m.cli),
		newColumn(m.ctx, textType, m.cli),
	}

	m.cols[credType].list.Title = "Credentials"
	m.cols[cardType].list.Title = "Bank cards"
	m.cols[textType].list.Title = "Notes"

	return nil
}

func NewMainBoard(ctx context.Context, cli *CLI) *MainBoard {
	help := help.New()
	help.ShowAll = true
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return &MainBoard{ctx: ctx, help: help, list: defaultList, cli: cli}
}

func (m *MainBoard) Init() tea.Cmd {
	return nil
}

func (m *MainBoard) updateColumns() error {

	resp, err := m.cli.Client.Resource.GetAllResources(m.ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	for _, res := range resp.GetResources() {
		switch res.GetType() {
		case string(models.Text):
			item := NewResourceItem(textType, int(res.Id), res.Description)
			m.cols[textType].Set(APPEND, item)
		case string(models.Card):
			item := NewResourceItem(cardType, int(res.Id), res.Description)
			m.cols[cardType].Set(APPEND, item)
		case string(models.Password):
			item := NewResourceItem(credType, int(res.Id), res.Description)
			m.cols[credType].Set(APPEND, item)
		}
	}
	return nil
}

func (m *MainBoard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		if !m.loaded {
			if err := m.updateColumns(); err != nil {
				if e, ok := status.FromError(err); ok {
					return NewErrorModel(e.Err()), nil
				}
			}
		}
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
		case key.Matches(msg, keys.Left):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getPrev()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Right):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext()
			m.cols[m.focused].Focus()
		}
	case *textForm:
		m.cols[textType].Set(APPEND, msg.createResource())
		return m, nil
	case *cardForm:
		m.cols[cardType].Set(APPEND, msg.createResource())
		return m, nil
	case *credForm:
		m.cols[credType].Set(APPEND, msg.createResource())
		return m, nil
	case ModelError:
		return msg.Update(nil)
	}

	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(column); ok {
		m.cols[m.focused] = res.(column)
	} else {
		return res, cmd
	}
	return m, cmd
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
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(&keys), m.help.ShortHelpView([]key.Binding{keys.New, keys.Enter, keys.Delete}))
}
