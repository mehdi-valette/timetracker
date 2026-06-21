package command

import "github.com/mehdi-valette/timetracker/internal/entity"

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
