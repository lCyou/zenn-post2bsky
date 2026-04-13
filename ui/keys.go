package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Post key.Binding
	Quit key.Binding
}

var DefaultKeyMap = KeyMap{
	Post: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "Post"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
}
