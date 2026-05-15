package command

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/entity"
)

func Run() {
	textinput := textinput.New()
	textinput.Focus()

	task := entity.CreateTask(0, "my task", entity.CreateDate())
	timeRange := entity.CreateTimeRange(0, 0, entity.CreateDate())
	task.SetTimeRange(timeRange)

	timeRange.Start()

	tea.NewProgram(model{
		input:       textinput,
		currentTask: task,
		clock:       Clock{},
	}).Run()
}

type model struct {
	input       textinput.Model
	currentTask entity.Tasker
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
			fmt.Print(m.input.Value())
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	case TickMsg:
		cmd := m.clock.Tick()
		return m, cmd
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

	return tea.NewView(m.input.View() + "\n" + taskInfo)
}
