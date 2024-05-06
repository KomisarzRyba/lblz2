package keymap

import "github.com/charmbracelet/bubbles/key"

type TableKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Filter key.Binding
	Quit   key.Binding
	Help   key.Binding
}

func (tkm TableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{tkm.Select, tkm.Filter, tkm.Quit, tkm.Help}
}

func (tkm TableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{tkm.Select, tkm.Filter},
		{tkm.Up, tkm.Down},
		{tkm.Help, tkm.Quit},
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
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

type DetailKeyMap struct {
	Print key.Binding
	Back  key.Binding
}

func (dkm DetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{dkm.Print, dkm.Back}
}

func (dkm DetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{dkm.Print, dkm.Back},
	}
}

func (dkm DetailKeyMap) WithReprint() DetailKeyMap {
	dkm.Print.SetKeys("r")
	dkm.Print.SetHelp("r", "reprint")
	return dkm
}

func NewDetailKeymap() DetailKeyMap {
	return DetailKeyMap{
		Print: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "print"),
		),
		Back: key.NewBinding(
			key.WithKeys("q", "esc", "backspace"),
			key.WithHelp("q/esc", "back"),
		),
	}
}
