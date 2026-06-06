package command

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/entity"
)

func Run() (returnModel tea.Model, returnErr error) {
	textinput := textinput.New()
	textinput.Focus()

	interpreter, _ := CreateTaskInterpreter()

	return tea.NewProgram(model{
		interpreter: interpreter,
		input:       textinput,
		currentTask: nil,
		information: "",
		clock:       Clock{},
	}).Run()
}

type model struct {
	interpreter Interpreter
	input       textinput.Model
	currentTask entity.Tasker
	information string
	error       error
	clock       Ticker
}

var _ tea.Model = model{}

func (m model) Init() tea.Cmd {
	return m.clock.Tick()
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyPressMsg:
		switch msg.Keystroke() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			// TODO: change the interpreter so it sends a new command instead of the value to show
			// this would allow to exit the program on Quit and change the current task
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
	case TaskCreatedMsg:
		m.error = nil
		m.currentTask = msg.task
		m.information = fmt.Sprintf("Created the task \"%s\"", msg.task.GetName())
		m.input = textinput.New()
		return m, m.input.Focus()
	case TaskListedMsg:
		m.error = nil
		m.information = msg.taskList
		m.input = textinput.New()
		return m, m.input.Focus()
	}

	return m, nil
}

func (m model) View() tea.View {
	taskInfo := "Current task: none"
	if m.currentTask != nil {
		duration, durationErr := m.currentTask.Duration()

		if durationErr != nil {
			panic(durationErr)
		}

		hours := duration / 3600
		minutes := (duration - (hours * 3600)) / 60
		seconds := duration - hours*3600 - minutes*60

		taskInfo = fmt.Sprintf("Current task (%d:%d:%d): %s", hours, minutes, seconds, m.currentTask.GetName())
	}

	text := m.input.View()

	if m.error != nil {
		text += "\n" + m.error.Error()
	}

	text += "\n" + taskInfo
	text += "\n" + m.information

	return tea.NewView(text)
}
