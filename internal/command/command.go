package command

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/command/alltasks"
	"github.com/mehdi-valette/timetracker/internal/command/common"
)

func Run(databasePath string) (returnModel tea.Model, returnErr error) {
	textinput := textinput.New()
	textinput.Focus()

	interpreter, _ := alltasks.CreateTaskInterpreter(databasePath)

	program := tea.NewProgram(
		alltasks.AllTasksModel{
			Interpreter: interpreter,
			Input:       textinput,
			Details:     viewport.Model{},
			Clock:       &common.Clock{},
		},
	)

	return program.Run()
}
