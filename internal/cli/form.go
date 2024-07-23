package cli

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sebasttiano/Owl/internal/models"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"google.golang.org/grpc/status"
)

var (
	ErrRenderTemplate = errors.New("error to render info template")
	ErrOutputType     = errors.New("unknown resource type")
)
var (
	//go:embed templates/card
	cardInfo string
	//go:embed templates/pass
	passInfo string
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
	description   textinput.Model
	content       textarea.Model
	cancelled     bool
	cli           *CLI
	help          help.Model
	resID         int
}

func newModel(ctx context.Context, cli *CLI, resType resType) tea.Model {
	switch resType {
	case 2:
		return newTextModel(ctx, cli)
	case 1:
		return newCardModel(ctx, cli)
	case 0:
		return newCredModel(ctx, cli)
	}

	return nil
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
	case ModelError:
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
				[]key.Binding{keys.Quit, keys.Save, keys.Back},
			),
		),
	)
}

func (f *textForm) createResource() ResourceItem {
	return NewResourceItem(textType, f.resID, f.description.Value())
}

type outputForm struct {
	ctx           context.Context
	width, height int
	cancelled     bool
	cli           *CLI
	help          help.Model
	resID         int
	col           *column
	content       string
	title         string
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

func (o *outputForm) getContent() tea.Model {

	request := &pb.GetResourceRequest{Id: int32(o.resID)}
	resp, _ := o.cli.Client.Resource.GetResource(o.ctx, request)
	res := resp.Resource
	switch res.GetType() {
	case string(models.Text):
		o.content = resp.Resource.GetContent()
	case string(models.Card):
		var card models.CardCreds
		if err := json.Unmarshal([]byte(res.GetContent()), &card); err != nil {
			return NewErrorModel(err)
		}

		tmpl, err := template.New("card").Parse(cardInfo)
		if err != nil {
			return NewErrorModel(ErrRenderTemplate)
		}
		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, card); err != nil {
			return NewErrorModel(ErrRenderTemplate)
		}
		o.content = buf.String()

	case string(models.Password):
		var creds models.Creds
		if err := json.Unmarshal([]byte(res.GetContent()), &creds); err != nil {
			return NewErrorModel(err)
		}

		tmpl, err := template.New("creds").Parse(passInfo)
		if err != nil {
			return NewErrorModel(ErrRenderTemplate)
		}
		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, creds); err != nil {
			return NewErrorModel(ErrRenderTemplate)
		}
		o.content = buf.String()
	default:
		return NewErrorModel(ErrOutputType)
	}
	o.title = resp.Resource.GetDescription()
	return nil
}

func (o *outputForm) View() string {
	return form(o.width, o.height, o.title, o.content)
}

type cardForm struct {
	ctx           context.Context
	width, height int
	description   textinput.Model
	ccn           textinput.Model
	exp           textinput.Model
	cvv           textinput.Model
	holder        textinput.Model
	cancelled     bool
	cli           *CLI
	help          help.Model
	resID         int
}

func newCardModel(ctx context.Context, cli *CLI) *cardForm {
	m := cardForm{
		ctx:         ctx,
		cli:         cli,
		description: textinput.New(),
		ccn:         textinput.New(),
		exp:         textinput.New(),
		cvv:         textinput.New(),
		holder:      textinput.New(),
		help:        help.New(),
	}
	m.description.CharLimit = 32
	m.description.Placeholder = "type your description"

	m.ccn.CharLimit = 16 + 3
	m.ccn.Placeholder = "4505 **** **** 1234"
	m.ccn.Prompt = ""
	m.ccn.Validate = func(s string) error {
		if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
			return fmt.Errorf("CCN is invalid")
		}
		if len(s)%5 == 0 && s[len(s)-1] != ' ' {
			return fmt.Errorf("CCN must separate groups with spaces")
		}
		c := strings.ReplaceAll(s, " ", "")
		_, err := strconv.ParseInt(c, 10, 64)
		return err
	}
	m.ccn.Focus()

	m.exp.CharLimit = 5
	m.exp.Placeholder = "MM/YY"
	m.exp.Prompt = ""
	m.exp.Validate = func(s string) error {
		e := strings.ReplaceAll(s, "/", "")
		_, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return fmt.Errorf("EXP is invalid")
		}
		if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
			return fmt.Errorf("EXP is invalid")
		}
		return nil
	}

	m.cvv.CharLimit = 3
	m.cvv.EchoMode = textinput.EchoPassword
	m.cvv.Placeholder = "123"
	m.cvv.Prompt = ""
	m.cvv.Validate = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	m.holder.CharLimit = 64
	m.holder.Placeholder = "CARD HOLDER"
	m.holder.Prompt = ""

	m.description.Focus()
	return &m
}

// Init implements tea.Model.
func (c *cardForm) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (c *cardForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		descriptionCmd tea.Cmd
		cnnCmd         tea.Cmd
		expCmd         tea.Cmd
		cvvCmd         tea.Cmd
		holderCmd      tea.Cmd
	)
	c.description, descriptionCmd = c.description.Update(msg)
	c.ccn, cnnCmd = c.ccn.Update(msg)
	c.exp, expCmd = c.exp.Update(msg)
	c.cvv, cvvCmd = c.cvv.Update(msg)
	c.holder, holderCmd = c.holder.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			switch {
			case c.description.Focused():
				if len(c.description.Value()) > 0 && c.description.Err == nil {
					c.description.Blur()
					c.ccn.Focus()
				}
			case c.ccn.Focused():
				if len(c.ccn.Value()) > 0 && c.ccn.Err == nil {
					c.ccn.Blur()
					c.exp.Focus()
				}
			case c.exp.Focused():
				if len(c.exp.Value()) > 0 && c.exp.Err == nil {
					c.exp.Blur()
					c.cvv.Focus()
				}
			case c.cvv.Focused():
				if len(c.cvv.Value()) > 0 && c.cvv.Err == nil {
					c.cvv.Blur()
					c.holder.Focus()
				}
			case c.holder.Focused():
				c.holder.Blur()
			}
		case key.Matches(msg, keys.Back):
			return mainBoard.Update(nil)
		case key.Matches(msg, keys.Quit):
			c.cancelled = true
			return c, tea.Quit
		case key.Matches(msg, keys.Save):
			if len(c.description.Value()) > 0 && len(c.ccn.Value()) > 0 && len(c.exp.Value()) > 0 && len(c.cvv.Value()) > 0 && len(c.holder.Value()) > 0 {
				card := &models.CardCreds{
					Description: c.description.Value(),
					CCN:         c.ccn.Value(),
					EXP:         c.exp.Value(),
					CVV:         c.cvv.View(),
					Holder:      c.holder.Value(),
				}
				return c, c.saveCardToServer(card)
			}
		}
	case *cardForm:
		return mainBoard.Update(c)
	case ModelError:
		return mainBoard.Update(msg)
	}
	return c, tea.Batch(cnnCmd, expCmd, cvvCmd, holderCmd, descriptionCmd)
}

func (c *cardForm) saveCardToServer(card *models.CardCreds) tea.Cmd {
	return func() tea.Msg {
		content, err := json.Marshal(card)
		if err != nil {
			return NewErrorModel(err)
		}
		request := pb.SetResourceRequest{Resource: &pb.ResourceMsg{Content: string(content), Description: card.Description, Type: string(models.Card)}}
		resp, err := c.cli.Client.Resource.SetResource(c.ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				return NewErrorModel(e.Err())
			}
		}
		c.resID = int(resp.Resource.GetId())
		return c
	}
}

// View implements tea.Model.
func (c *cardForm) View() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#002288"))
	return form(
		c.width, c.height,
		"Card",
		lipgloss.JoinVertical(
			lipgloss.Center,
			c.description.View(),
			strings.Repeat(" ", 64),
			lipgloss.JoinVertical(
				lipgloss.Left,
				hintStyle.Render("Card Number"),
				c.ccn.View(),
				" ",
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					lipgloss.JoinVertical(
						lipgloss.Left,
						hintStyle.Render("Expiration Date"),
						c.exp.View(),
					),
					"       ",
					lipgloss.JoinVertical(
						lipgloss.Left,
						hintStyle.Render("CVV"),
						c.cvv.View(),
					),
				),
				" ",
				hintStyle.Render("Card Holder"),
				c.holder.View(),
				" ",
				help.New().ShortHelpView(
					[]key.Binding{keys.Quit, keys.Enter, keys.Save, keys.Back},
				),
			),
		))
}

func (c *cardForm) createResource() ResourceItem {
	return NewResourceItem(cardType, c.resID, c.description.Value())
}

type credForm struct {
	ctx           context.Context
	width, height int
	description   textinput.Model
	username      textinput.Model
	password      textinput.Model
	cancelled     bool
	cli           *CLI
	help          help.Model
	resID         int
}

func newCredModel(ctx context.Context, cli *CLI) *credForm {
	m := credForm{
		ctx:         ctx,
		cli:         cli,
		description: textinput.New(),
		username:    textinput.New(),
		password:    textinput.New(),
		help:        help.New(),
	}
	m.description.CharLimit = 32
	m.description.Placeholder = "type your description"
	m.description.Focus()

	m.username.CharLimit = 32
	m.username.Prompt = "Username: "
	m.username.Placeholder = "type username..."

	m.password.CharLimit = 32
	m.password.Prompt = "Password: "
	m.password.Placeholder = "type password..."
	m.password.EchoMode = textinput.EchoPassword
	return &m
}

func (c *credForm) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (c *credForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		descriptionCmd tea.Cmd
		usernameCmd    tea.Cmd
		passwordCmd    tea.Cmd
	)
	c.description, descriptionCmd = c.description.Update(msg)
	c.username, usernameCmd = c.username.Update(msg)
	c.password, passwordCmd = c.password.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.help.Width = c.width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			switch {
			case c.description.Focused():
				if len(c.description.Value()) > 0 && c.description.Err == nil {
					c.description.Blur()
					c.username.Focus()
				}
			case c.username.Focused():
				if len(c.username.Value()) > 0 && c.username.Err == nil {
					c.username.Blur()
					c.password.Focus()
				}
			case c.password.Focused():
				if len(c.password.Value()) > 0 && c.password.Err == nil {
					c.password.Blur()
				}
			}
		case key.Matches(msg, keys.Back):
			return mainBoard.Update(nil)
		case key.Matches(msg, keys.Quit):
			c.cancelled = true
			return c, tea.Quit
		case key.Matches(msg, keys.Save):
			if len(c.description.Value()) > 0 && len(c.username.Value()) > 0 && len(c.password.Value()) > 0 {
				creds := &models.Creds{
					Description: c.description.Value(),
					Username:    c.username.Value(),
					Password:    c.password.Value(),
				}
				return c, c.saveCredToServer(creds)
			}
		}
	case *credForm:
		return mainBoard.Update(c)
	case ModelError:
		return mainBoard.Update(msg)
	}
	return c, tea.Batch(descriptionCmd, usernameCmd, passwordCmd)
}

// View implements tea.Model.
func (c *credForm) View() string {
	return form(
		c.width, c.height,
		"Credentials",
		lipgloss.JoinVertical(
			lipgloss.Left,
			c.description.View(),
			strings.Repeat(" ", 64),
			c.username.View(),
			strings.Repeat(" ", 64),
			c.password.View(),
			strings.Repeat(" ", 64),
			c.help.ShortHelpView(
				[]key.Binding{keys.Quit, keys.Save, keys.Back},
			),
		),
	)
}

func (c *credForm) createResource() ResourceItem {
	return NewResourceItem(credType, c.resID, c.description.Value())
}

func (c *credForm) saveCredToServer(creds *models.Creds) tea.Cmd {
	return func() tea.Msg {
		content, err := json.Marshal(creds)
		if err != nil {
			return NewErrorModel(err)
		}
		request := pb.SetResourceRequest{Resource: &pb.ResourceMsg{Content: string(content), Description: creds.Description, Type: string(models.Password)}}
		resp, err := c.cli.Client.Resource.SetResource(c.ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				return NewErrorModel(e.Err())
			}
		}
		c.resID = int(resp.Resource.GetId())
		return c
	}
}
