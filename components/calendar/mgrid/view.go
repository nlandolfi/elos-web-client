package mgrid

import (
	"time"

	"github.com/nlandolfi/spin/apps/cal"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme *ui.Theme
	Time  *time.Time
}

func View(s *State) *browser.Node {
	return s.Theme.Text("hellow mgrid")
}

func Grid(s *GridSpec) *browser.Node {
	if s == nil {
		panic("Grid: GridSpec is nil")
	}

	days := cal.DaysInMonthGridOf(s.Time)

	m := int(len(days) / 7)

	if s.PadTo6Weeks && m < 6 {
		for j := time.Duration(0); j < 7; j++ {
			days = append(days, days[m*7-1].Add((j+1)*24*time.Hour))
		}
		m += 1
	}

	dayLabels := make([]*browser.Node, 7)
	for i := 0; i < 7; i++ {
		dayLabels[i] = s.Head(days, i)
	}

	var rows []*browser.Node = []*browser.Node{
		ui.HStack(dayLabels...).FlexGrow("1").FlexBasis("0px").AlignItemsCenter(),
	}

	for i := 0; i < m; i++ {
		var cols []*browser.Node

		for j := 0; j < 7; j++ {
			d := days[i*7+j]
			cols = append(cols, s.Cell(days, i, j, d))
		}

		rows = append(rows, ui.HStack(cols...).FlexGrow("1").FlexBasis("0px").AlignItemsCenter())
	}

	return ui.VStack(rows...).FlexGrow("1")
}

type GridSpec struct {
	Time        time.Time
	Head        func(days []time.Time, i int) *browser.Node
	Cell        func(days []time.Time, i, j int, t time.Time) *browser.Node
	PadTo6Weeks bool
}
