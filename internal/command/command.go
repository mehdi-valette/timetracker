package command

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

func Run(databasePath string) (returnModel tea.Model, returnErr error) {
	textinput := textinput.New()
	textinput.Focus()

	interpreter, _ := CreateTaskInterpreter(databasePath)

	return tea.NewProgram(timeTrackerModel{
		interpreter: interpreter,
		input:       textinput,
		information: "",
		clock:       &Clock{},
	}).Run()
}

type timeTrackerModel struct {
	interpreter Interpreter
	input       textinput.Model
	information string
	error       error
	clock       Ticker
}

var _ tea.Model = timeTrackerModel{}

func (m timeTrackerModel) Init() tea.Cmd {
	return m.clock.Tick()
}

func (m timeTrackerModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyPressMsg:
		switch msg.Keystroke() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, m.interpreter.Interpret(m.input.Value())
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	case TickMsg:
		cmd := m.clock.Tick()
		return m, cmd
	case ErrorMsg:
		m.error = msg.error
		m.information = ""
	case TaskCreatedMsg:
		m.error = nil
		m.information = fmt.Sprintf("Created the task (%d) %s", msg.task.GetId(), msg.task.GetName())
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskListedMsg:
		m.error = nil
		m.information = msg.taskList
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskStartedMsg:
		m.error = nil
		m.information = fmt.Sprintf("Started the task (%d) %s", msg.task.GetId(), msg.task.GetName())
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskStoppedMsg:
		m.error = nil
		m.information = fmt.Sprintf("Stopped the task (%d) %s", msg.task.GetId(), msg.task.GetName())
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskDeletedMsg:
		m.error = nil
		m.information = fmt.Sprintf("Deleted the task (%d) %s", msg.task.GetId(), msg.task.GetName())
		m.input = textinput.New()
		return m, m.input.Focus()
	}

	return m, nil
}

func (m timeTrackerModel) View() tea.View {
	taskInfo := "Current task: none"
	currentTask := m.interpreter.GetCurrentTask()

	if currentTask != nil {
		duration, durationErr := currentTask.Duration()

		if durationErr != nil {
			panic(durationErr)
		}

		taskInfo = fmt.Sprintf("Current task (%d | %s) %s", currentTask.GetId(), duration.ToString(), currentTask.GetName())
	}

	text := m.input.View()

	text += "\n" + taskInfo
	if m.error != nil {
		text += "\n" + m.error.Error()
	}
	text += "\n--------------------\n" + m.information

	return tea.NewView(text)
}
