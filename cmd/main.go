package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shopspring/decimal"

	"github.com/gurgeous/vectro/internal"
)

//nolint:recvcheck // bubbletea required Update
type Model struct {
	args Args
	// underlying calculator for math
	c *internal.Calculator
	// window size
	width  int
	height int
	// error to display in red, or "say" message in green
	err string
	say string
	// text input, and is it visible?
	input        textinput.Model
	inputVisible bool
	// vhs mode (demo.tape)
	vhs       bool
	vhsTyping bool
	vhsBanner string
}

func InitModel() Model {
	return InitModelWithArgs(Args{})
}

func InitModelWithArgs(args Args) Model {
	// only ParseArgs if not testing, since `go test` passed -test.paniconexit0
	m := Model{
		args: args,
		c:    internal.NewCalculator(),
		input: func() textinput.Model {
			input := textinput.New()
			input.Focus()
			input.Placeholder = "enter number..."
			input.Width = 20
			input.Cursor.Style = internal.CursorStyle
			return input
		}(),
		vhs: os.Getenv("VHS") != "",
	}

	if m.vhs {
		m.args.noInit = true
	}

	return m
}

func (m Model) Init() tea.Cmd {
	if !m.args.noInit {
		Load(m.c)
	}
	return textinput.Blink
}

//
// Update
//

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// always clear msg when the user hits a key
		m.err, m.say = "", ""

		if m.vhs && m.vhsUpdate(msg) {
			return m, nil
		}

		// paste?
		if msg.Paste {
			m.paste(string(msg.Runes))
			return m, nil
		}

		// quit?
		if slices.Contains(QuitKeys, msg.String()) {
			if !m.args.noInit {
				Save(m.c)
			}
			return m, tea.Quit
		}

		// some other key
		var err error
		cmd, err = m.onKey(msg)
		if err != nil {
			m.err = err.Error()
			m.inputVisible = false
			m.input.Reset()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	default:
		// for blinking
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}

//
// handle a keypress
//

var (
	// these keys show the numeric input
	NumberKeys = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "."}
	// these keys quit
	QuitKeys = []string{"q", "ctrl+c", "ctrl+q"}
)

func (m *Model) onKey(msg tea.KeyMsg) (tea.Cmd, error) {
	var cmd tea.Cmd

	key := msg.String()
	if command, ok := internal.CommandsByKey[key]; ok {
		return cmd, m.run(command.Name)
	}

	// non-input keys
	if !m.inputVisible {
		if key == "backspace" {
			return cmd, m.run(internal.DROP)
		}
		if key == "enter" {
			return cmd, m.run(internal.DUP)
		}
		if slices.Contains(NumberKeys, key) {
			m.inputVisible = true
		}
	}

	// input keys
	if m.inputVisible {
		if key == "enter" {
			return cmd, m.enter(true)
		}
		m.input, cmd = m.input.Update(msg)
	}

	return cmd, nil
}

//
// View
//
// Note that with responsive sizing various boxes can be hidden.
//
// 111 333  1=stack
// 111 333  2=history
// 222 333  3=help
// 222 333
// 4444444  4=status
//
//

func (m Model) View() string {
	// get screen width, bail early if not available yet
	w, h := m.width, m.height
	if w == 0 || h == 0 {
		return ""
	}

	//
	// boxes
	//

	boxAll := internal.NewBox(w, h)
	boxMain, box4 := boxAll.CutBottom(1)
	boxLeft, box3 := boxMain.Cols()
	if box3.GetWidth() < 40 {
		// too narrow? hide help
		boxLeft, box3 = boxMain, internal.NewBox(0, 0)
	}
	// border + padding.vert + stack + text + border
	stackHeight := internal.StackStyle.GetVerticalPadding() + 1 + internal.StackSize + 1 + 1
	box1, box2 := boxLeft.CutTop(stackHeight)
	if box2.GetHeight() < 5 {
		// too short? hide history
		box1, box2 = boxLeft, internal.NewBox(0, 0)
	}

	if box1.GetWidth() < 20 || box1.GetHeight() < stackHeight {
		// too cramped?
		style := boxAll.Apply(internal.CrampedStyle)
		return style.Render("vectro is feeling cramped, make your terminal bigger!")
	}

	//
	// boxes => styles
	//

	style1 := box1.Apply(internal.StackStyle)
	style2 := box2.Apply(internal.PaneStyle)
	style3 := box3.Apply(internal.PaneStyle)
	style4 := box4.Apply(internal.StatusStyle)

	//
	// render
	//

	str1 := RenderPane(style1, m.title(), m.stack(style1))
	var str2 string
	if !m.vhs {
		str2 = RenderPane(style2, "history", m.history(style2))
	} else {
		str2 = RenderPane(internal.BannerStyle.Inherit(style2), "demo", m.vhsBanner)
	}
	str3 := RenderPane(style3, "keys", m.help(style3))
	str4 := style4.Render(m.status(style4))

	var left string
	if str2 != "" {
		left = lipgloss.JoinVertical(0, str1, str2)
	} else {
		left = str1
	}
	return lipgloss.JoinVertical(0, lipgloss.JoinHorizontal(0, left, str3), str4)
}

// handle vhs stuff
func (m *Model) vhsUpdate(msg tea.KeyMsg) bool {
	key := msg.String()

	if !m.vhsTyping {
		if key == "[" {
			// enter vhs typing mode
			m.vhsTyping = true
			if len(m.vhsBanner) > 0 {
				m.vhsBanner += "\n"
			}
			return true
		}
		if key == "ctrl+e" {
			// clear screen
			m.vhsBanner = ""
			return true
		}
		return false
	}

	if key == "]" {
		// exit vhs typing mode
		m.vhsTyping = false
		return true
	}

	m.vhsBanner += key
	return true
}

// handle enter key (or the programmatic equivalent)
func (m *Model) enter(explicit bool) error {
	if m.input.Value() != "" {
		val, err := decimal.NewFromString(m.input.Value())
		if err != nil {
			return errors.New("invalid number")
		}
		m.c.Enter(val, explicit)
	}
	m.inputVisible = false
	m.input.Reset()
	return nil
}

func (m *Model) inputNeg() error {
	s := m.input.Value()
	if len(s) == 0 {
		return fmt.Errorf("%s: %s", internal.NEG, "too few arguments")
	}
	switch {
	case strings.HasPrefix(s, "-"):
		s = "+" + s[1:]
	case strings.HasPrefix(s, "+"):
		s = "-" + s[1:]
	default:
		s = "-" + s
	}
	m.input.SetValue(s)
	m.input.CursorEnd()
	return nil
}

func (m *Model) paste(str string) {
	re := regexp.MustCompile(`[^\d.+-]`)
	paste := re.ReplaceAllString(str, "")
	if paste != "" {
		if m.inputVisible {
			m.input.SetValue(m.input.Value() + paste)
		} else {
			m.inputVisible = true
			m.input.SetValue(paste)
		}
		m.input.CursorEnd()
	}
}

func (m *Model) run(name string) error {
	// implicit ENTER, maybe
	if m.inputVisible {
		if name == internal.NEG {
			return m.inputNeg()
		}
		if name == internal.UNDO {
			m.say = "undo"
			m.inputVisible = false
			m.input.Reset()
			return nil
		}
		if err := m.enter(false); err != nil {
			return err
		}
	}
	if err := m.c.Run(name); err != nil {
		return fmt.Errorf("%s: %s", name, err.Error())
	}
	if name == internal.YANK {
		m.say = "yanked to clipboard"
	}
	if name == internal.UNDO {
		m.say = "undo"
	}

	return nil
}

//
// title/stack/history/help/status rendering
//

func (m Model) title() string {
	if m.err != "" {
		return internal.ErrorStyle.Render(m.err)
	}
	if m.say != "" {
		return internal.SayStyle.Render(m.say)
	}

	var s string
	var chars = "Vectro"
	for ii, rune := range chars {
		s += internal.TitleStyles[ii].Render(string(rune))
	}
	return s
}

func (m Model) stack(style lipgloss.Style) string {
	stack := internal.MapWithIndex(m.c.GetDisplay(), func(ii int, str string) string {
		array := strings.Split(str, ":")
		return internal.IndexStyle.Render(array[0]+":") + internal.GradientStyles[ii].Render(array[1])
	})
	if m.inputVisible {
		stack = internal.Push(stack, " "+m.input.View())
	}
	return strings.Join(internal.ClipLines(stack, style), "\n")
}

func (m Model) history(style lipgloss.Style) string {
	history := m.c.History()
	history = internal.Reverse(internal.ClipLines(internal.Reverse(history), style))
	return strings.Join(history, "\n")
}

//go:embed help.txt
var StaticHelpText string

func (m Model) help(style lipgloss.Style) string {
	// clip/wrap
	w := style.GetWidth() - style.GetHorizontalPadding()
	h := style.GetHeight() - style.GetVerticalPadding()
	plain := lipgloss.NewStyle().
		Width(w).
		Height(h).
		MaxWidth(w).
		MaxHeight(h).
		Render(StaticHelpText)
	return internal.StyleBetweenStars(plain, internal.HelpKeyStyle)
}

func (m Model) status(style lipgloss.Style) string {
	w := style.GetWidth() - style.GetHorizontalPadding()
	if w < 60 {
		return "vectro"
	}
	return "https://github.com/gurgeous/vectro"
}

//
// main
//

// set at build time by goreleaser
var (
	date    = "today"
	version = "development"
)

func main() {
	args := ParseArgs(os.Args[1:])
	p := tea.NewProgram(InitModelWithArgs(args), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
