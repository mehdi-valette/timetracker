package states

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/command/alltasks"
	"github.com/mehdi-valette/timetracker/internal/command/common"
)

var slicePopErr = errors.New("cannot pop from an empty array")

type States struct {
	models struct {
		AllTasks alltasks.AllTasksModel
	}
	Clock        common.Ticker
	LastModels   []*tea.Model
	CurrentModel tea.Model
}

var _ tea.Model = States{}

// Init implements [tea.Model].
func (s States) Init() tea.Cmd {
	return s.Clock.Tick()
}

// Update implements [tea.Model].
func (s States) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return s.CurrentModel.Update(message)
}

// View implements [tea.Model].
func (s States) View() tea.View {
	return s.CurrentModel.View()
}

func slicePop[S []E, E any](slice S) (E, S, error) {
	if len(slice) == 0 {
		var zeroValue E
		return zeroValue, nil, slicePopErr
	}

	lastIndex := len(slice) - 1

	return slice[lastIndex], slice[:lastIndex], nil
}
