package command

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/command/alltasks"
	"github.com/mehdi-valette/timetracker/internal/command/common"
)

type Tmp struct{}

func (tmp Tmp) Init() tea.Cmd {
	return func() tea.Msg { return "" }
}

func (tmp Tmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return Tmp{}, tea.Cmd(nil)
}

func (tmp Tmp) View() tea.View {
	return tea.View{Content: ""}
}

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
