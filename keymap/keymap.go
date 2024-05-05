package keymap

import "github.com/charmbracelet/bubbles/key"

type TableKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
	Help   key.Binding
}

func (tkm TableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{tkm.Select, tkm.Quit, tkm.Help}
}

func (tkm TableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{tkm.Select, tkm.Quit},
		{tkm.Up, tkm.Down},
		{tkm.Help},
	}
}

func NewTableKeyMap() TableKeyMap {
	return TableKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

type ItemKeyMap struct {
	Print key.Binding
	Back  key.Binding
}

func (ikm ItemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{ikm.Print, ikm.Back}
}

// TODO: check if required
// func (ikm ItemKeyMap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{ikm.Print, ikm.Back},
// 	}
// }
