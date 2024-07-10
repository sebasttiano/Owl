package cli

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sebasttiano/Owl/internal/models"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"google.golang.org/grpc/status"
	"strings"
)

var formStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#bb0000")).
	Padding(1, 2, 1, 2)

func form(width, height int, title, content string) string {
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			formStyle.Render(content),
		),
	)
}

type textForm struct {
	ctx           context.Context
	width, height int
	index         int
	description   textinput.Model
	content       textarea.Model
	cancelled     bool
	cli           *CLI
	help          help.Model
	resID         int
}

func newTextModel(ctx context.Context, cli *CLI) *textForm {
	m := textForm{
		ctx:         ctx,
		cli:         cli,
		description: textinput.New(),
		content:     textarea.New(),
		help:        help.New(),
	}
	m.description.CharLimit = 32
	m.description.Placeholder = "type your description"
	m.content.ShowLineNumbers = false
	m.content.MaxWidth = 64
	m.content.CharLimit = 1024
	m.content.Placeholder = "type your note..."
	m.content.SetHeight(12)
	m.content.SetWidth(64)

	m.description.Focus()
	return &m
}

// Init implements tea.Model.
func (f *textForm) Init() tea.Cmd {
	return textarea.Blink
}

// Update implements tea.Model.
func (f *textForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var contentCmd tea.Cmd
	var descriptionCmd tea.Cmd
	f.content, contentCmd = f.content.Update(msg)
	f.description, descriptionCmd = f.description.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
		f.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			f.cancelled = false
			return f, tea.Quit
		case key.Matches(msg, keys.Back):
			return mainBoard.Update(nil)
		case key.Matches(msg, keys.Enter):
			switch {
			case f.description.Focused():
				if len(f.description.Value()) < 1 {
					return f, textinput.Blink
				}
				f.description.Blur()
				f.content.Focus()
			case f.content.Focused():
				if len(f.content.Value()) < 1 {
					return f, textinput.Blink
				}
				f.content.Blur()
			}
		case key.Matches(msg, keys.Save):
			if len(f.description.Value()) > 0 && len(f.content.Value()) > 0 {
				return f, f.saveTextToServer(f.description.Value(), f.content.Value())
			}
		}
	case *textForm:
		return mainBoard.Update(f)
	case ErrorModel:
		return mainBoard.Update(msg)
	}

	return f, tea.Batch(descriptionCmd, contentCmd)
}

func (f *textForm) saveTextToServer(description, content string) tea.Cmd {
	return func() tea.Msg {
		request := pb.SetResourceRequest{Resource: &pb.ResourceMsg{Content: content, Description: description, Type: string(models.Text)}}
		resp, err := f.cli.Client.Resource.SetResource(f.ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				return NewErrorModel(e.Err())
			}
		}

		f.resID = int(resp.Resource.GetId())
		return f
	}
}

// View implements tea.Model.
func (f *textForm) View() string {
	return form(
		f.width, f.height,
		"Resource Note",
		lipgloss.JoinVertical(
			lipgloss.Left,
			f.description.View(),
			strings.Repeat(" ", 64),
			f.content.View(),
			lipgloss.PlaceHorizontal(
				lipgloss.Width(f.content.View()),
				lipgloss.Right,
				fmt.Sprintf("%d/%d", f.content.Length(), f.content.CharLimit),
			),
			" ",
			f.help.ShortHelpView(
				[]key.Binding{
					key.NewBinding(
						key.WithKeys("esc"),
						key.WithHelp("[esc]", "cancel"),
					),
					key.NewBinding(
						key.WithKeys("ctrl+s"),
						key.WithHelp("[ctrl+s]", "save"),
					),
				},
			),
		),
	)
}

func (f *textForm) createResource() ResourceItem {
	item := ResourceItem{resType: textType, description: f.description.Value(), resID: f.resID}
	item.title = item.MakeTitle()
	return item
}

type outputForm struct {
	ctx           context.Context
	width, height int
	index         int
	cancelled     bool
	cli           *CLI
	help          help.Model
	resID         int
	col           *column
	content       string
}

func newOutputModel(ctx context.Context, cli *CLI, col *column, id int) *outputForm {
	o := &outputForm{
		col:   col,
		ctx:   ctx,
		cli:   cli,
		resID: id,
	}
	return o
}

func (o *outputForm) Init() tea.Cmd {
	return nil
}

func (o *outputForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		o.width = msg.Width
		o.height = msg.Height
		o.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			o.cancelled = true
			return o, tea.Quit
		case key.Matches(msg, keys.Back):
			return mainBoard.Update(nil)
		}
	}
	return o, nil
}

func (o *outputForm) getContent() {
	request := &pb.GetResourceRequest{Id: int32(o.resID)}
	resp, _ := o.cli.Client.Resource.GetResource(o.ctx, request)
	o.content = resp.Resource.GetContent()
}

func (o *outputForm) View() string {
	return form(o.width, o.height, o.col.list.Title, o.content)
	//return form(
	//	o.width, o.height,
	//	"Register to Owl",
	//	lipgloss.JoinVertical(
	//		lipgloss.Left,
	//		s.username.View(),
	//		strings.Repeat(" ", 64),
	//		s.password.View(),
	//		strings.Repeat(" ", 64),
	//		s.repeatPassword.View(),
	//	),
	//)
}
