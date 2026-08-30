package alltasks

import (
	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/entity"
)

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{error: err}
	}
}

type ErrorMsg struct {
	error error
}

type HelpMsg struct {
	help string
}

type TaskStartedMsg struct {
	task entity.Tasker
}

type TaskStoppedMsg struct {
	task entity.Tasker
}

type TaskCreatedMsg struct {
	task entity.Tasker
}

type TaskBeganMsg struct {
	task entity.Tasker
}

type TaskListedMsg struct {
	taskList string
}

type TaskDeletedMsg struct {
	task entity.Tasker
}

type TaskRenamed struct {
	task entity.Tasker
}
