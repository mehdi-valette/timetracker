package command

import tea "charm.land/bubbletea/v2"

type ErrorMsg struct {
	error error
}

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{error: err}
	}
}
