package day

import (
	"time"

	"github.com/nlandolfi/elos/web-client/components/calendar/ltr"
	"github.com/nlandolfi/spin/apps/cal"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme      *ui.Theme        `json:"-"`
	Time       *time.Time       `json:"-"`
	EventItems []*cal.EventItem `json:"-"` // this is a pointer

	DispatchEditEvent func(*cal.EventItem) `json:"-"`

	hoveredEventID string
}

func (s *State) Handle(e browser.Event) {
	switch e := e.(type) {
	case EventDecrementDay:
		*s.Time = s.Time.Add(-24 * time.Hour)
	case EventSetTime:
		*s.Time = e.Time
	case EventIncrementDay:
		*s.Time = s.Time.Add(24 * time.Hour)
	case EventEventHoverStart:
		s.hoveredEventID = e.ID
	case EventEventHoverLeave:
		if s.hoveredEventID == e.ID {
			s.hoveredEventID = ""
		}
	}
}

type EventDecrementDay struct{}
type EventSetTime struct{ time.Time }
type EventIncrementDay struct{}
type EventEventHoverStart struct{ ID string }
type EventEventHoverLeave struct{ ID string }

func (s *State) SetTheme(t *ui.Theme) {
	s.Theme = t
}

func View(s *State) *browser.Node {
	var es cal.EventItems

	for _, e := range s.EventItems {
		if e.ShouldDisplayOnDay(*s.Time) {
			es = append(es, e)
		}
	}

	cards := make([]*browser.Node, len(es))

	for i, e := range es {
		cards[i] = dayCard(s, e)
	}

	return ui.VStack(
		ui.HStack(
			s.Theme.Text(s.Time.Format("Monday January 2, 2006")).
				FontSizeEM(1.5),
			ui.OnlyIf(cal.SameDay(*s.Time, time.Now()), // todo pull current time from state?
				func() *browser.Node {
					return s.Theme.Text(s.Time.Format("(today)")).
						FontSizeEM(1).
						MarginLeftPX(10)
				},
			),
			ui.Spacer(),
			ltr.View(&ltr.State{
				Theme:        s.Theme,
				OnClickPrev:  browser.Dispatcher(EventDecrementDay{}),
				OnClickToday: browser.Dispatcher(EventSetTime{time.Now()}),
				OnClickNext:  browser.Dispatcher(EventIncrementDay{}),
			}),
		).AlignItemsCenter().FlexWrap(browser.FlexWrapWrap),
		ui.If(len(cards) > 0,
			func() *browser.Node {
				return ui.VStack(cards...)
			},
			func() *browser.Node {

				return s.Theme.Text("No events on this day...")
			},
		),
	)
}

func dayCard(s *State, item *cal.EventItem) *browser.Node {
	return s.Theme.Card(
		ui.HStack(
			s.Theme.Text(item.Name).
				FontSize(browser.Size{Value: 1.2, Unit: browser.UnitEM}),
			ui.Spacer(),
			s.Theme.Text(item.Time.Format("2 Jan 2006")),
		),
		ui.OnlyIf(item.HourSpecified,
			func() *browser.Node {
				return s.Theme.Text(item.Time.Format("3:04 PM"))
			},
		),
		ui.If(len(item.Details) == 0,
			func() *browser.Node {
				return s.Theme.Text("No details...").Color("lightgray").FontSizeEM(0.5) // TODO: pull color from theme
			},
			func() *browser.Node {
				return s.Theme.Text(item.Details)
			},
		),
		ui.If(item.Recurs,
			func() *browser.Node {
				return s.Theme.Text(item.RecurString())
			},
			func() *browser.Node {
				return s.Theme.Text("No recurrence...").Color("lightgray").FontSizeEM(0.5) // TODO: pull color from theme
			},
		),
		ui.HStack(
			ui.Spacer(),
			s.Theme.Button("Edit").OnClick(func(_ dom.Event) {
				if s.DispatchEditEvent != nil {
					go s.DispatchEditEvent(item)
				}
			}),
		),
	).PaddingPX(10)
}
