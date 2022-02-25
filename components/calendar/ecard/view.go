package ecard

import (
	"github.com/nlandolfi/spin/apps/cal"
	"github.com/nlandolfi/spin/web/browser"
	"github.com/nlandolfi/spin/web/browser/ui"
)

type State struct {
	Theme *ui.Theme
	Item  *cal.EventItem
	// below is not implemented
	RecurrenceNumber int // if set, the display will make it look like this events time is this reccurence
}

func View(s *State) *browser.Node {
	item := *s.Item

	return s.Theme.Card(
		ui.VStack(
			ui.HStack(
				s.Theme.Text(item.Name),
				ui.If(
					!item.HourSpecified,
					s.Theme.Text(item.DateString()),
					s.Theme.Textf("%s @ %s", item.DateString(), item.TimeString()),
				),
			),
		).JustifyContentSpaceBetween(),
		s.Theme.Text(item.Details),
		ui.OnlyIf(item.Recurs,
			ui.HStack(
				s.Theme.Text("This event recurs:").Padding(3),
				s.Theme.Text(item.RecurString()).Padding(3),
				ui.OnlyIf(item.UntilSpecified,
					s.Theme.Textf("Until: %s", item.Until.Format("2 Jan 2006")),
				),
			),
		),
	).Padding(10)
}
