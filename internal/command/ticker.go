package command

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type TickMsg struct{}

type Ticker interface {
	Tick() tea.Cmd
}

type Clock struct {
	lastTick time.Time
}

var _ Ticker = Clock{}

func (t Clock) Tick() tea.Cmd {
	ellapsedTime := time.Since(t.lastTick)
	remainingTime := 1000 - ellapsedTime.Milliseconds()

	t.lastTick = time.Now()

	return tea.Tick(time.Millisecond*time.Duration(remainingTime), func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}
