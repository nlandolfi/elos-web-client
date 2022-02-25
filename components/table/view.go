package table

import (
	"time"

	"github.com/nlandolfi/spin/apps/cal"
	"github.com/nlandolfi/spin/web/browser"
	"github.com/nlandolfi/spin/web/browser/ui"
)

type State struct {
	Theme      *ui.Theme        `json:"-"`
	EventItems []*cal.EventItem `json:"-"` // this is a pointer
	// TODO maybe in the future have a pointer to current time.
}

func (s *State) SetTheme(t *ui.Theme) {
	s.Theme = t
}

func View(s *State) *browser.Node {
	t := time.Now()

	var es cal.EventItems

	for _, e := range s.EventItems {
		// TODO: handle recurs and duration events, this is just a start
		if e.Time.After(t) {
			es = append(es, e)
		}
	}

	views := make([]*browser.Node, len(es))
	for i, e := range es {
		views[i] = row(s, e.CSVEntry())
	}

	return ui.VStack(
		row(s, cal.CSVHeaderEntry).FontWeight("bold"),
		ui.VStack(views...),
	)
}

func rowe(s *State, e *cal.EventItem) *browser.Node {
	var views []*browser.Node
	for _, entry := range e.CSVEntry() {
		views = append(views, s.Theme.Text(entry))
	}
	return ui.HStack(views...)
}

func row(s *State, entry []string) *browser.Node {
	var views []*browser.Node
	for _, col := range entry {
		views = append(
			views,
			s.Theme.Text(col).
				FlexBasis("0px").
				FlexGrow("1").
				BorderRight(border(s.Theme)),
		)
	}
	return ui.HStack(views...)
}

func border(th *ui.Theme) browser.Border {
	return browser.Border{
		Color: th.TextColor,
		Width: browser.Size{Value: 1, Unit: browser.UnitPX},
		Type:  browser.BorderSolid,
	}
}
