package command

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func Run(databasePath string) (returnModel tea.Model, returnErr error) {
	textinput := textinput.New()
	textinput.Focus()

	interpreter, _ := CreateTaskInterpreter(databasePath)

	program := tea.NewProgram(timeTrackerModel{
		interpreter: interpreter,
		input:       textinput,
		details:     viewport.Model{},
		clock:       &Clock{},
	})

	return program.Run()
}

type timeTrackerModel struct {
	interpreter Interpreter
	input       textinput.Model
	details     viewport.Model
	error       error
	clock       Ticker
}

var _ tea.Model = timeTrackerModel{}

func (m timeTrackerModel) Init() tea.Cmd {
	return m.clock.Tick()
}

func (m timeTrackerModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := message.(type) {
	case tea.KeyPressMsg:
		switch msg.Keystroke() {
		case "esc", "ctrl+c":
			cmds = append(cmds, m.interpreter.Interpret("quit"))
		case "ctrl+v":
			cmds = append(cmds, func() tea.Msg { return tea.PasteMsg{} })
		case "ctrl+l":
			m.input = textinput.New()
			cmds = append(cmds, m.input.Focus())
		case "enter":
			cmds = append(cmds, m.interpreter.Interpret(m.input.Value()))
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.PasteMsg:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(textinput.Paste())
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		cmds = append(cmds, m.interpreter.ListTasks())
		m.details.SetWidth(msg.Width)
		m.details.SetHeight(msg.Height - lipgloss.Height(m.HeaderView()))
	case HelpMsg:
		m.details.SetContent(msg.help)
	case TickMsg:
		cmds = append(cmds, m.clock.Tick())
	case ErrorMsg:
		m.error = msg.error
	case TaskCreatedMsg:
		m.error = nil
		m.input = textinput.New()
		cmds = append(cmds, m.input.Focus())
		cmds = append(cmds, m.interpreter.ListTasks())
	case TaskListedMsg:
		m.error = nil
		m.details.SetContent(msg.taskList)
		m.input = textinput.New()
		cmds = append(cmds, m.input.Focus())
	case TaskStartedMsg:
		m.error = nil
		m.input = textinput.New()
		cmds = append(cmds, m.input.Focus())
		cmds = append(cmds, m.interpreter.ListTasks())
	case TaskStoppedMsg:
		m.error = nil
		m.input = textinput.New()
		cmds = append(cmds, m.input.Focus())
		cmds = append(cmds, m.interpreter.ListTasks())
	case TaskDeletedMsg:
		m.error = nil
		m.input = textinput.New()
		cmds = append(cmds, m.input.Focus())
		cmds = append(cmds, m.interpreter.ListTasks())
	}

	var cmd tea.Cmd
	m.details, cmd = m.details.Update(message)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m timeTrackerModel) HeaderView() string {
	currentTaskInfo := "Current task: none"

	if currentTask := m.interpreter.GetCurrentTask(); currentTask != nil {
		currentTaskInfo = fmt.Sprintf("Current task: %s", currentTask.String())
	}

	err := ""

	if m.error != nil {
		err = m.error.Error()
	} else if m.input.Err != nil {
		err = m.input.Err.Error()
	}

	return lipgloss.NewStyle().Render(fmt.Sprintf("%s\n%s\n--------------------\n%s\n\n", m.input.View(), err, currentTaskInfo))
}

func (m timeTrackerModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	v.SetContent(fmt.Sprintf("%s%s", m.HeaderView(), m.details.View()))

	return v
}
