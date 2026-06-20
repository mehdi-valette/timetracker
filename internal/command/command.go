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

	program := tea.NewProgram(timeTrackerModel{
		interpreter: interpreter,
		input:       textinput,
		information: "",
		clock:       &Clock{},
	})

	return program.Run()
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
			return m, m.interpreter.Interpret("quit")
		case "ctrl+v":
			return m, func() tea.Msg { return tea.PasteMsg{} }
		case "ctrl+l":
			m.input = textinput.New()
			m.input.Focus()
		case "enter":
			return m, m.interpreter.Interpret(m.input.Value())
		default:
			m.input, _ = m.input.Update(msg)
		}
	case tea.PasteMsg:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(textinput.Paste())
		return m, cmd
	case TickMsg:
		cmd := m.clock.Tick()
		return m, cmd
	case ErrorMsg:
		m.error = msg.error
		m.information = ""
	case TaskCreatedMsg:
		m.error = nil
		m.information = fmt.Sprintf("Created the task %s", msg.task.StringShort())
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskListedMsg:
		m.error = nil
		m.information = msg.taskList
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskStartedMsg:
		m.error = nil
		m.information = fmt.Sprintf("Started the task %s", msg.task.StringShort())
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskStoppedMsg:
		m.error = nil
		m.information = fmt.Sprintf("Stopped the task %s", msg.task.StringShort())
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskDeletedMsg:
		m.error = nil
		m.information = fmt.Sprintf("Deleted the task %s", msg.task.StringShort())
		m.input = textinput.New()
		return m, m.input.Focus()
	}

	return m, nil
}

func (m timeTrackerModel) View() tea.View {
	taskInfo := "Current task: none"
	currentTask := m.interpreter.GetCurrentTask()

	if currentTask != nil {
		taskInfo = fmt.Sprintf("Current task %s", currentTask.String())
	}

	info := m.information

	if m.error != nil {
		info = m.error.Error()
	} else if m.input.Err != nil {
		info = m.input.Err.Error()
	}

	return tea.NewView(fmt.Sprintf("%s\n%s\n-------------------\n%s", m.input.View(), taskInfo, info))
}
