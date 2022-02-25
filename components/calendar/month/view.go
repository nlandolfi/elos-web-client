package month

import (
	"time"

	"github.com/nlandolfi/elos/web-client/components/calendar/ltr"
	"github.com/nlandolfi/spin/apps/cal"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme          *ui.Theme        `json:"-"`
	Time           *time.Time       `json:"-"`
	SelectedKey    *string          `json:"-"`
	EventItems     []*cal.EventItem `json:"-"` // this is a pointer
	hoveredDay     time.Time
	hoveredEventID string

	InspectEvent func(e *cal.EventItem) `json:"-"`
}

func (s *State) Handle(e browser.Event) {
	switch e := e.(type) {
	case EventDayClick:
		*s.SelectedKey = "day"
		*s.Time = e.Time
	case EventDayHoverStart:
		s.hoveredDay = e.Time
	case EventDayHoverLeave:
		if s.hoveredDay.Equal(e.Time) {
			s.hoveredDay = time.Time{}
		}
	case EventDecrementMonth:
		m, y := s.Time.Month(), s.Time.Year()

		if m == time.January {
			y -= 1
			m = time.December
		} else {
			m -= 1
		}

		// the three is a hack to avoid time zone nonsense for now
		*s.Time = time.Date(y, m, 3, 0, 0, 0, 0, s.Time.Location())
	case EventSetTime:
		*s.Time = e.Time
	case EventIncrementMonth:
		m, y := s.Time.Month(), s.Time.Year()

		if m == time.December {
			y += 1
			m = time.January
		} else {
			m += 1
		}

		// the three is a hack to avoid time zone nonsense for now
		*s.Time = time.Date(y, m, 2, 0, 0, 0, 0, s.Time.Location())
	case EventEventHoverStart:
		s.hoveredEventID = e.ID
	case EventEventHoverLeave:
		if s.hoveredEventID == e.ID {
			s.hoveredEventID = ""
		}
	case EventInspectEvent:
		s.InspectEvent(e.EventItem)
	}
}

type EventDayClick struct{ time.Time }
type EventDayHoverStart struct{ time.Time }
type EventDayHoverLeave struct{ time.Time }
type EventDecrementMonth struct{}
type EventSetTime struct{ time.Time }
type EventIncrementMonth struct{}
type EventEventHoverStart struct{ ID string }
type EventEventHoverLeave struct{ ID string }
type EventInspectEvent struct{ *cal.EventItem }

func (s *State) SetTheme(t *ui.Theme) {
	s.Theme = t
}

func View(s *State) *browser.Node {
	return view(s)
}

func view(s *State) *browser.Node {
	return ui.VStack(
		ui.HStack(
			s.Theme.Text(s.Time.Format("January, 2006")).
				FontSizeEM(1.5).
				MarginLeftPX(10),
			/*
				ui.OnlyIf(s.Time.Month() == time.Now().Month(),
					func() *browser.Node {
						return s.Theme.Text(s.Time.Format("(this month)")).
							FontSizeEM(1).
							MarginLeftPX(10)
					}),
			*/
			ui.Spacer(),
			ltr.View(&ltr.State{
				Theme:        s.Theme,
				OnClickPrev:  browser.Dispatcher(EventDecrementMonth{}),
				OnClickToday: browser.Dispatcher(EventSetTime{time.Now()}),
				OnClickNext:  browser.Dispatcher(EventIncrementMonth{}),
			}),
		).AlignItemsCenter(),
		grid(s),
	).FlexGrow("1")
}

func grid(s *State) *browser.Node {
	days := cal.DaysInMonthGridOf(*s.Time)

	m := int(len(days) / 7)

	dayLabels := make([]*browser.Node, 7)
	for i := 0; i < 7; i++ {
		dayLabels[i] = s.Theme.Text(days[i].Format("Mon")).
			TextAlignRight().
			MarginRightPX(5).
			FlexGrow("1").
			FlexBasis("0px")
	}
	var rows []*browser.Node = []*browser.Node{
		ui.HStack(dayLabels...),
	}

	for i := 0; i < m; i++ {
		var cols []*browser.Node

		for j := 0; j < 7; j++ {
			d := days[i*7+j]
			cols = append(cols, gridSquare(s, d, i == 0, j == 0, i == m-1, j == 7-1))
		}

		rows = append(rows, ui.HStack(cols...).FlexGrow("1"))
	}

	//	log.Printf("%s is a %s; month has %d days", day, day.Weekday(), len(days))

	//	return ui.VStack(s.Theme.Text("NONSENSE"))

	return ui.VStack(rows...).FlexGrow("1")
}

func gridSquare(s *State, t time.Time, rowStart, colStart, rowEnd, colEnd bool) *browser.Node {
	var es cal.EventItems
	for _, e := range s.EventItems {
		if e.ShouldDisplayOnDay(t) {
			es = append(es, e)
		}
	}
	views := make([]*browser.Node, len(es))
	for i, e := range es {
		views[i] = lineView(s, e).CursorPointer()
	}

	return ui.VStack(
		dayLineView(s, t),
		ui.VStack(views...),
	).
		FlexGrow("1").
		FlexBasis("0px").
		MinWidth(browser.Size{Value: 60, Unit: browser.UnitPX}).
		MinHeight(browser.Size{Value: 100, Unit: browser.UnitPX}).
		BorderTop(border(s.Theme)).
		BorderLeft(border(s.Theme)).
		OnlyIf(rowEnd, func(n *browser.Node) *browser.Node {
			return n.BorderBottom(border(s.Theme))
		}).
		OnlyIf(colEnd, func(n *browser.Node) *browser.Node {
			return n.BorderRight(border(s.Theme))
		})
}

func dayLineView(s *State, t time.Time) *browser.Node {
	lcache := t.String()
	return ui.HStack(
		ui.Spacer(),
		ui.If(t.Day() == 1,
			func() *browser.Node { return s.Theme.Textf("%s %d", t.Format("Jan"), t.Day()) },
			func() *browser.Node { return s.Theme.Textf("%d", t.Day()) },
		).
			TextAlignRight().
			OnlyIf(s.Time.Month() != t.Month(),
				func(n *browser.Node) *browser.Node {
					return n.Color("gray")
								}).
			OnlyIf(cal.SameDay(t, time.Now()), // TODO get time from state
				func(n *browser.Node) *browser.Node {
					return n.Color("red") // TODO use theme
				}).
			OnlyIf(cal.SameDay(s.hoveredDay, t),
				func(n *browser.Node) *browser.Node {
					return n.Background("lightgray").BorderRadiusPX(40) // TODO use theme
				}).
			PaddingLeftPX(8).
			PaddingRightPX(8).
			MarginRightPX(5).
			MarginTopPX(5).
			FontSizeEM(0.8),
	).
		OnMouseOverCached(lcache, browser.Dispatcher(EventDayHoverStart{t})).
		OnMouseOutCached(lcache, browser.Dispatcher(EventDayHoverLeave{t})).
		OnClickCached(lcache, browser.Dispatcher(EventDayClick{t})).
		Pointer()
}

func border(th *ui.Theme) browser.Border {
	return browser.Border{
		Color: "lightgray", // TODO: use theme?
		Width: browser.Size{Value: 1, Unit: browser.UnitPX},
		Type:  browser.BorderSolid,
	}
}

func lineView(s *State, e *cal.EventItem) *browser.Node {
	return s.Theme.Text(e.Name).OverflowHidden().
		FontSizeEM(1).
		MaxHeight(browser.Size{Value: 1.2, Unit: browser.UnitEM}).
		Padding(browser.Size{Value: 2, Unit: browser.UnitPX}).
		Background("lightgray").
		Color("black").
		MarginTopPX(2).
		OnMouseOverCached(e.ID, browser.Dispatcher(EventEventHoverStart{e.ID})).
		OnMouseOutCached(e.ID, browser.Dispatcher(EventEventHoverLeave{e.ID})).
		OnlyIf(e.ID == s.hoveredEventID, func(n *browser.Node) *browser.Node {
			return n.BoxShadow(itemLiftedShadow)
		}).
		OnClickCached(e.ID, browser.Dispatcher(EventInspectEvent{e})).
		MarginBottomPX(1)
}

var itemLiftedShadow = browser.BoxShadow{
	HOffset: browser.Size{},
	VOffset: browser.Size{Value: 4, Unit: browser.UnitPX},
	Blur:    browser.Size{Value: 8, Unit: browser.UnitPX},
	Spread:  browser.Size{Value: 0, Unit: browser.UnitPX},
	Color:   "rgba(0, 0, 0, 0.1)",
}
