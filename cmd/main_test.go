package main

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	m := InitModel()

	// try typing something
	assert.False(t, m.inputVisible)
	m, _ = testUpdate(m, testKeyMsg("1"))
	m, _ = testUpdate(m, testKeyMsg("."))
	m, _ = testUpdate(m, testKeyMsg("2"))
	assert.True(t, m.inputVisible)
	// backspace
	m, _ = testUpdate(m, testKeyMsg("9"))
	m, _ = testUpdate(m, testKeyMsg("backspace"))
	// enter
	m, _ = testUpdate(m, testKeyMsg("enter"))
	assert.Equal(t, 1.2, m.c.PopFloat64())
	assert.False(t, m.inputVisible)

	// implicit enter
	m, _ = testUpdate(m, testKeyMsg("9"))
	assert.True(t, m.inputVisible)
	m, _ = testUpdate(m, testKeyMsg("@"))
	assert.Equal(t, 3, m.c.PopInt())

	// quit
	_, msg := testUpdate(m, testKeyMsg("q"))
	assert.Equal(t, reflect.TypeOf(msg()), reflect.TypeOf(tea.QuitMsg{}))
}

func TestMainCommands(t *testing.T) {
	m := InitModel()

	// DIV
	m.c.PushInt(27, 3)
	m, _ = testUpdate(m, testKeyMsg("/"))
	assert.Equal(t, 9, m.c.PopInt())
	assert.True(t, m.c.Empty())

	// NEG
	m.c.PushInt(123)
	m, _ = testUpdate(m, testKeyMsg("n"))
	assert.Equal(t, -123, m.c.PopInt())
	assert.True(t, m.c.Empty())
}

func TestMainMisc(t *testing.T) {
	m := InitModel()

	// enter
	m.c.PushInt(123)
	m, _ = testUpdate(m, testKeyMsg("enter"))
	assert.Equal(t, 123, m.c.PopInt())
	assert.Equal(t, 123, m.c.PopInt())

	// drop
	m, _ = testUpdate(m, testKeyMsg("backspace"))
	m, _ = testUpdate(m, testKeyMsg("backspace"))
	assert.True(t, m.c.Empty())

	// yank (commented out to avoid disturbin the clipboard while testing0
	// m.c.PushInt(456)
	// m, _ = testUpdate(m, testKeyMsg("y"))
	// clip, _ := clipboard.ReadAll()
	// assert.Equal(t, "456", clip)
}

func TestNeg(t *testing.T) {
	m := InitModel()
	m, _ = testUpdate(m, testKeyMsg("1"))
	m, _ = testUpdate(m, testKeyMsg("."))
	m, _ = testUpdate(m, testKeyMsg("2"))
	m, _ = testUpdate(m, testKeyMsg("n"))
	assert.Equal(t, "-1.2", m.input.Value())
	m, _ = testUpdate(m, testKeyMsg("n"))
	assert.Equal(t, "+1.2", m.input.Value())
}

func TestRendering(t *testing.T) {
	m := InitModel()

	// make sure everything is visible
	m.width, m.height = 80, 40
	view := ansi.Strip(m.View())
	assert.Contains(t, view, "1:")
	assert.Contains(t, view, "history")
	assert.Contains(t, view, "keys")
	assert.Contains(t, view, "github")

	// narrow? no help
	m.width, m.height = 20, 40
	view = ansi.Strip(m.View())
	assert.Contains(t, view, "1:")
	assert.Contains(t, view, "history")
	assert.NotContains(t, view, "keys")

	// short? no help
	m.width, m.height = 80, 15
	view = ansi.Strip(m.View())
	assert.Contains(t, view, "1:")
	assert.NotContains(t, view, "history")
	assert.Contains(t, view, "keys")

	// both? no keys or help
	m.width, m.height = 20, 15
	view = ansi.Strip(m.View())
	assert.Contains(t, view, "1:")
	assert.NotContains(t, view, "history")
	assert.NotContains(t, view, "keys")

	// cramped
	m.width, m.height = 15, 10
	view = ansi.Strip(m.View())
	assert.Contains(t, view, "cramped")
}

//
// helpers
//

func testKeyMsg(key string) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

func testUpdate(m Model, msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.Update(msg)
	return model.(Model), cmd
}
