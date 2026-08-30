package alltasks

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mehdi-valette/timetracker/internal/command/common"
)

type AllTasksModel struct {
	Interpreter Interpreter
	Input       textinput.Model
	Details     viewport.Model
	Error       error
	Clock       common.Ticker
}

var _ tea.Model = AllTasksModel{}

func (m AllTasksModel) Init() tea.Cmd {
	return m.Clock.Tick()
}

func (m AllTasksModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := message.(type) {
	case tea.KeyPressMsg:
		switch msg.Keystroke() {
		case "esc", "ctrl+c":
			cmds = append(cmds, m.Interpreter.Interpret("quit"))
		case "ctrl+v":
			cmds = append(cmds, func() tea.Msg { return tea.PasteMsg{} })
		case "ctrl+l":
			m.Input = textinput.New()
			cmds = append(cmds, m.Input.Focus())
		case "enter":
			cmds = append(cmds, m.Interpreter.Interpret(m.Input.Value()))
		default:
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.PasteMsg:
		m.Input, cmd = m.Input.Update(textinput.Paste())
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		cmds = append(cmds, m.Interpreter.ListTasks())
		m.Details.SetWidth(msg.Width)
		m.Details.SetHeight(msg.Height - lipgloss.Height(m.HeaderView()))
	case HelpMsg:
		m.Details.SetContent(msg.help)
	case common.TickMsg:
		cmds = append(cmds, m.Clock.Tick())
	case ErrorMsg:
		m.Error = msg.error
	case TaskListedMsg:
		m.Error = nil
		m.Details.SetContent(msg.taskList)
		m.Input = textinput.New()
		cmds = append(cmds, m.Input.Focus())
	case TaskDeletedMsg, TaskRenamed, TaskStoppedMsg, TaskStartedMsg, TaskCreatedMsg:
		m.Error = nil
		m.Input = textinput.New()
		cmds = append(cmds, m.Input.Focus())
		cmds = append(cmds, m.Interpreter.ListTasks())
	}

	m.Details, cmd = m.Details.Update(message)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m AllTasksModel) HeaderView() string {
	currentTaskInfo := "Current task: none"

	if currentTask := m.Interpreter.GetCurrentTask(); currentTask != nil {
		currentTaskInfo = fmt.Sprintf("Current task: %s", currentTask.String())
	}

	err := ""

	if m.Error != nil {
		err = m.Error.Error()
	} else if m.Input.Err != nil {
		err = m.Input.Err.Error()
	}

	return lipgloss.NewStyle().Render(fmt.Sprintf("%s\n%s\n--------------------\n%s\n\n", m.Input.View(), err, currentTaskInfo))
}

func (m AllTasksModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	v.SetContent(fmt.Sprintf("%s%s", m.HeaderView(), m.Details.View()))

	return v
}
