package cli

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	ErrCanceledByUser       = errors.New("authentication canceled by user")
	ErrUnknownModelType     = errors.New("unknown model type")
	ErrRegisterPassMismatch = errors.New("passwords are not the same")
)

var signBoard *SignBoard

func (c *CLI) GetUserCreds(ctx context.Context, tlsCreds credentials.TransportCredentials) (string, string, error) {

	signBoard = NewSignBoard(c, c.cfg.Info.Banner, tlsCreds)
	m, err := tea.NewProgram(
		signBoard,
		tea.WithAltScreen(),
		tea.WithContext(ctx)).Run()

	if err != nil {
		return "", "", err
	}

	model, ok := m.(signInModel)
	if !ok {
		return "", "", ErrUnknownModelType
	}

	if model.cancelled {
		return "", "", ErrCanceledByUser
	}
	return model.username.Value(), model.password.Value(), nil

}

type SignBoard struct {
	cancelled     bool
	list          list.Model
	help          help.Model
	width, height int
	loaded        bool
	cli           *CLI
	banner        string
}

func NewSignBoard(c *CLI, banner string, tlsCreds credentials.TransportCredentials) *SignBoard {
	help := help.New()
	help.ShowAll = true
	signList := list.New([]list.Item{
		newSignInModel("SIGN IN", ""),
		newSignUpModel("SIGN UP", "", c, tlsCreds),
	}, list.NewDefaultDelegate(), 0, 0)
	signList.SetShowHelp(false)
	signList.Title = "Do you have an account?"

	return &SignBoard{help: help, list: signList, banner: banner}
}

func (s *SignBoard) Init() tea.Cmd {
	return nil
}

func (s *SignBoard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.height = msg.Height
		s.width = msg.Width
		s.loaded = true
		s.list.SetSize(s.width/margin, s.height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			s.cancelled = true
			return s, tea.Quit
		case key.Matches(msg, keys.Enter):
			switch option := s.list.SelectedItem().(type) {
			case signInModel:
				return option.Update(nil)
			case signUpModel:
				return option.Update(nil)
			}
		}

	case signUpModel:
		s.list, cmd = s.list.Update(msg)
		return s, cmd
	}
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s *SignBoard) View() string {
	if s.cancelled {
		return ""
	}
	if !s.loaded {
		return "LOADING!"
	}
	return form(s.width, s.height, s.banner,
		lipgloss.JoinVertical(
			lipgloss.Center,
			s.list.View(),
			s.help.View(&keys)))
}

type signUpModel struct {
	cli            *CLI
	cancelled      bool
	width, height  int
	username       textinput.Model
	password       textinput.Model
	repeatPassword textinput.Model
	item           Item
	tls            credentials.TransportCredentials
}

// Init implements tea.Model.
func (s signUpModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (s signUpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		usernameCmd   tea.Cmd
		passwordCmd   tea.Cmd
		repeatPassCmd tea.Cmd
	)
	s.username, usernameCmd = s.username.Update(msg)
	s.password, passwordCmd = s.password.Update(msg)
	s.repeatPassword, repeatPassCmd = s.repeatPassword.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.height = msg.Height
		s.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			switch {
			case s.username.Focused():
				if len(s.username.Value()) < 1 {
					return s, textinput.Blink
				}
				s.username.Blur()
				s.password.Focus()
			case s.password.Focused():
				if len(s.password.Value()) < 1 {
					return s, textinput.Blink
				}
				s.password.Blur()
				s.repeatPassword.Focus()
			case s.repeatPassword.Focused():
				if len(s.repeatPassword.Value()) < 1 {
					return s, textinput.Blink
				}
				s.repeatPassword.Blur()
				if s.password.Value() != s.repeatPassword.Value() {
					return NewErrorModel(ErrRegisterPassMismatch), nil
				}

				return s, s.makeRegistration
			}
		case key.Matches(msg, keys.Quit):
			s.cancelled = true
			return s, tea.Quit
		}
	case ModelError:
		return msg.Update(nil)
	case *SignBoard:
		msg.width = s.width
		msg.height = s.height
		msg.loaded = true
		return msg.Update(nil)
	}

	return s, tea.Batch(usernameCmd, passwordCmd, repeatPassCmd)
}

func (s signUpModel) makeRegistration() tea.Msg {

	regConn, err := grpc.NewClient(
		s.cli.cfg.GetServerAddress(),
		grpc.WithTransportCredentials(s.tls))

	if err != nil {
		return NewErrorModel(err)
	}

	s.cli.Client, err = NewGRPCClient(regConn)
	if err != nil {
		return NewErrorModel(err)
	}

	request := pb.RegisterRequest{Name: s.username.Value(), Password: s.password.Value()}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.cli.Client.Auth.Register(ctx, &request)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return NewErrorModel(e.Err())
		}
	}
	return signBoard
}

// View implements tea.Model.
func (s signUpModel) View() string {
	return form(
		s.width, s.height,
		"Register to Owl",
		lipgloss.JoinVertical(
			lipgloss.Left,
			s.username.View(),
			strings.Repeat(" ", 64),
			s.password.View(),
			strings.Repeat(" ", 64),
			s.repeatPassword.View(),
		),
	)
}

// implement the list.Item interface
func (s signUpModel) FilterValue() string {
	return s.item.title
}

func (s signUpModel) Title() string {
	return s.item.title
}

func (s signUpModel) Description() string {
	return s.item.description
}

func newSignUpModel(title, description string, c *CLI, tlsCreds credentials.TransportCredentials) signUpModel {
	m := signUpModel{
		cancelled:      false,
		username:       textinput.New(),
		password:       textinput.New(),
		repeatPassword: textinput.New(),
		item: Item{
			title:       title,
			description: description,
		},
		cli: c,
		tls: tlsCreds,
	}
	m.username.CharLimit = 32
	m.username.Prompt = "Username: "
	m.username.Placeholder = "type your username..."

	m.password.CharLimit = 32
	m.password.Prompt = "Password: "
	m.password.EchoMode = textinput.EchoPassword
	m.password.Placeholder = "type your password..."

	m.repeatPassword.CharLimit = 32
	m.repeatPassword.Prompt = "Repeat Password: "
	m.repeatPassword.EchoMode = textinput.EchoPassword
	m.repeatPassword.Placeholder = "repeat your password again..."

	m.username.Focus()
	return m
}

type signInModel struct {
	item          Item
	cancelled     bool
	width, height int

	username textinput.Model
	password textinput.Model
}

// Init implements tea.Model.
func (s signInModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (s signInModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		usernameCmd tea.Cmd
		passwordCmd tea.Cmd
	)
	s.username, usernameCmd = s.username.Update(msg)
	s.password, passwordCmd = s.password.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.height = msg.Height
		s.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			switch {
			case s.username.Focused():
				if len(s.username.Value()) < 1 {
					return s, textinput.Blink
				}
				s.username.Blur()
				s.password.Focus()
			case s.password.Focused():
				if len(s.password.Value()) < 1 {
					return s, textinput.Blink
				}
				s.password.Blur()
				return s, tea.Quit
			}
		case key.Matches(msg, keys.Quit):
			s.cancelled = true
			return s, tea.Quit
		}
	}
	return s, tea.Batch(usernameCmd, passwordCmd)
}

// View implements tea.Model.
func (s signInModel) View() string {
	return form(
		s.width, s.height,
		"Login to Owl",
		lipgloss.JoinVertical(
			lipgloss.Left,
			s.username.View(),
			strings.Repeat(" ", 64),
			s.password.View(),
		),
	)
}

// FilterValue implement the list.Item interface
func (s signInModel) FilterValue() string {
	return s.item.title
}

func (s signInModel) Title() string {
	return s.item.title
}

func (s signInModel) Description() string {
	return s.item.description
}

func newSignInModel(title, description string) signInModel {
	m := signInModel{
		cancelled: false,
		username:  textinput.New(),
		password:  textinput.New(),
		item: Item{
			title:       title,
			description: description,
		},
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
