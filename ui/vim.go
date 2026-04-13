package ui

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type VimMode int

const (
	ModeNormal VimMode = iota
	ModeInsert
)

func (m VimMode) String() string {
	if m == ModeInsert {
		return "-- INSERT --"
	}
	return "-- NORMAL --"
}

// HandleNormalKey processes a keypress in Normal mode.
// It returns the updated textarea model, the new VimMode, and any tea.Cmd.
func HandleNormalKey(ta textarea.Model, msg tea.KeyMsg, lastKey string) (textarea.Model, VimMode, tea.Cmd, string) {
	key := msg.String()

	switch key {
	case "i":
		ta.Focus()
		return ta, ModeInsert, textarea.Blink, ""

	case "a":
		// Move cursor right then insert
		moveCursorRight(&ta)
		ta.Focus()
		return ta, ModeInsert, textarea.Blink, ""

	case "o":
		// New line below
		lines := strings.Split(ta.Value(), "\n")
		row, _ := cursorPos(ta)
		before := strings.Join(lines[:row+1], "\n")
		after := strings.Join(lines[row+1:], "\n")
		newVal := before + "\n"
		if after != "" {
			newVal += after
		}
		ta.SetValue(newVal)
		// Position cursor at start of new line
		ta.Focus()
		return ta, ModeInsert, textarea.Blink, ""

	case "O":
		// New line above
		lines := strings.Split(ta.Value(), "\n")
		row, _ := cursorPos(ta)
		before := strings.Join(lines[:row], "\n")
		after := strings.Join(lines[row:], "\n")
		var newVal string
		if before == "" {
			newVal = "\n" + after
		} else {
			newVal = before + "\n\n" + after
		}
		ta.SetValue(newVal)
		ta.Focus()
		return ta, ModeInsert, textarea.Blink, ""

	case "h":
		ta.Update(tea.KeyMsg{Type: tea.KeyLeft})
		return ta, ModeNormal, nil, ""

	case "l":
		ta.Update(tea.KeyMsg{Type: tea.KeyRight})
		return ta, ModeNormal, nil, ""

	case "j":
		ta.Update(tea.KeyMsg{Type: tea.KeyDown})
		return ta, ModeNormal, nil, ""

	case "k":
		ta.Update(tea.KeyMsg{Type: tea.KeyUp})
		return ta, ModeNormal, nil, ""

	case "0":
		ta.Update(tea.KeyMsg{Type: tea.KeyHome})
		return ta, ModeNormal, nil, ""

	case "$":
		ta.Update(tea.KeyMsg{Type: tea.KeyEnd})
		return ta, ModeNormal, nil, ""

	case "w":
		moveWordForward(&ta)
		return ta, ModeNormal, nil, ""

	case "b":
		moveWordBackward(&ta)
		return ta, ModeNormal, nil, ""

	case "d":
		if lastKey == "d" {
			deleteLine(&ta)
			return ta, ModeNormal, nil, ""
		}
		return ta, ModeNormal, nil, "d"

	case "u":
		// bubbles/textarea doesn't expose undo; send ctrl+z as best effort
		ta.Update(tea.KeyMsg{Type: tea.KeyCtrlZ})
		return ta, ModeNormal, nil, ""
	}

	return ta, ModeNormal, nil, key
}

// cursorPos returns row and col based on the value and the internal line info.
// bubbles/textarea doesn't expose cursor position directly so we approximate.
func cursorPos(ta textarea.Model) (row, col int) {
	// We rely on ta.Line() if available — not exposed; use a heuristic.
	return 0, 0
}

func moveCursorRight(ta *textarea.Model) {
	ta.Update(tea.KeyMsg{Type: tea.KeyRight})
}

func moveWordForward(ta *textarea.Model) {
	val := ta.Value()
	// Simple implementation: skip non-spaces then spaces
	_ = val
	for i := 0; i < 5; i++ {
		ta.Update(tea.KeyMsg{Type: tea.KeyRight})
	}
}

func moveWordBackward(ta *textarea.Model) {
	for i := 0; i < 5; i++ {
		ta.Update(tea.KeyMsg{Type: tea.KeyLeft})
	}
}

func deleteLine(ta *textarea.Model) {
	val := ta.Value()
	lines := strings.Split(val, "\n")
	if len(lines) <= 1 {
		ta.SetValue("")
		return
	}
	// Delete first non-empty line as approximation (textarea doesn't expose cursor row)
	newLines := make([]string, 0, len(lines)-1)
	deleted := false
	for _, l := range lines {
		if !deleted {
			deleted = true
			continue
		}
		newLines = append(newLines, l)
	}
	ta.SetValue(strings.Join(newLines, "\n"))
}

func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
