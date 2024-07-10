package cli

import (
	"context"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"google.golang.org/grpc/status"
)

const APPEND = -1

type column struct {
	ctx     context.Context
	focus   bool
	resType resType
	list    list.Model
	height  int
	width   int
	cli     *CLI
}

func (c *column) Focus() {
	c.focus = true
}

func (c *column) Blur() {
	c.focus = false
}

func (c *column) Focused() bool {
	return c.focus
}

func newColumn(ctx context.Context, resType resType, cli *CLI) column {
	var focus bool
	if resType == credType {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return column{ctx: ctx, cli: cli, focus: focus, resType: resType, list: defaultList}
}

// Init does initial setup for the column.
func (c column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/margin, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.New):
			t := newTextModel(c.ctx, c.cli)
			t.index = APPEND
			return t.Update(nil)
		case key.Matches(msg, keys.Enter):
			if len(c.list.VisibleItems()) != 0 {
				item := c.list.SelectedItem().(ResourceItem)
				o := newOutputModel(c.ctx, c.cli, &c, item.resID)
				o.getContent()
				return o.Update(nil)
			}
		case key.Matches(msg, keys.Delete):
			return c, c.DeleteCurrent()
			//case key.Matches(msg, keys.Enter):
			//	return c, c.MoveToNext()
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *column) DeleteCurrent() tea.Cmd {
	if len(c.list.VisibleItems()) > 0 {
		item := c.list.SelectedItem().(ResourceItem)
		request := &pb.DeleteTextRequest{Id: int32(item.resID)}
		_, err := c.cli.Client.Text.DeleteText(c.ctx, request)
		if err != nil {
			if _, ok := status.FromError(err); ok {
				return tea.Quit
			}
		}
		c.list.RemoveItem(c.list.Index())
	}
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)
	return cmd
}

func (c *column) Set(i int, r ResourceItem) tea.Cmd {
	if i != APPEND {
		return c.list.SetItem(i, r)
	}
	return c.list.InsertItem(APPEND, r)
}

func (c *column) setSize(width, height int) {
	c.width = width / margin
}

func (c *column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}
