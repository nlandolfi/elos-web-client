package week

import (
	"log"
	"time"

	"github.com/nlandolfi/elos/web-client/components/calendar/ltr"
	"github.com/nlandolfi/spin/apps/cal"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme       *ui.Theme  `json:"-"`
	Time        *time.Time `json:"-"`
	SelectedKey *string
	EventItems  []*cal.EventItem `json:"-"` // this is a pointer

	hoveredDay     time.Time
	hoveredEventID string
}

func (s *State) SetTheme(t *ui.Theme) {
	s.Theme = t
}

type EventEventClick struct {
	*cal.EventItem
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
	case EventDecrementWeek:
		*s.Time = s.Time.Add(-7 * 24 * time.Hour)
	case EventSetTime:
		*s.Time = e.Time
	case EventIncrementWeek:
		*s.Time = s.Time.Add(7 * 24 * time.Hour)
	case EventEventHoverStart:
		s.hoveredEventID = e.ID
	case EventEventHoverLeave:
		if s.hoveredEventID == e.ID {
			s.hoveredEventID = ""
		}
	}
}

type EventDayClick struct{ time.Time }
type EventDayHoverStart struct{ time.Time }
type EventDayHoverLeave struct{ time.Time }
type EventDecrementWeek struct{}
type EventSetTime struct{ time.Time }
type EventIncrementWeek struct{}
type EventEventHoverStart struct{ ID string }
type EventEventHoverLeave struct{ ID string }
type EventInspectEvent struct{ *cal.EventItem }

func View(s *State) *browser.Node {
	return view(s)
}

func view(s *State) *browser.Node {
	return ui.VStack(
		ui.HStack(
			s.Theme.Text(s.Time.Format("Week of January 2, 2006")).
				FontSizeEM(1.5).
				MarginLeftPX(10),
			ui.Spacer(),
			ltr.View(&ltr.State{
				Theme:        s.Theme,
				OnClickPrev:  browser.Dispatcher(EventDecrementWeek{}),
				OnClickToday: browser.Dispatcher(EventSetTime{time.Now()}),
				OnClickNext:  browser.Dispatcher(EventIncrementWeek{}),
			}),
		).AlignItemsCenter(),
		simplegrid(s),
	).FlexGrow("1")
}

func simplegrid(s *State) *browser.Node {

	dayCols := make([]*browser.Node, 7)

	for i, day := range cal.DaysInWeekOf(*s.Time) {
		lcache := day.String()
		var eviews []*browser.Node = []*browser.Node{
			ui.HStack(
				s.Theme.Button(day.Format("Mon 2")).
					OnClickCached(lcache, browser.Dispatcher(EventDayClick{day})),
			).BorderBottom(border(s.Theme)).JustifyContentCenter(),
		}
		for _, e := range s.EventItems {
			if e.ShouldDisplayOnDay(day) {
				eviews = append(eviews, lineView(s, e))
			}
		}

		dayCols[i] = ui.VStack(eviews...).
			FlexGrow("1").
			FlexBasis("0px").
			MarginLeftPX(1).
			MarginRightPX(1)
	}

	return ui.HStack(dayCols...)
}

func grid(s *State) *browser.Node {
	var days []time.Time

	for i := 0; i < int(s.Time.Weekday()); i++ {
		log.Printf("%d before", i)
		days = append(days, s.Time.Add(-time.Duration(int(s.Time.Weekday())-i)*24*time.Hour))
	}
	days = append(days, *s.Time)
	for i := s.Time.Weekday() + 1; i <= time.Saturday; i++ {
		log.Printf("%d after", i)
		days = append(days, s.Time.Add(time.Duration(i)*24*time.Hour))
	}

	var cols []*browser.Node

	cols = append(cols, ui.VStack(
		s.Theme.Text("").FlexGrow("1"),
		s.Theme.Text("all-day").FlexGrow("1"),
		s.Theme.Text("12:00 AM").FlexGrow("1"),
		s.Theme.Text(" 1:00 AM").FlexGrow("1"),
		s.Theme.Text(" 2:00 AM").FlexGrow("1"),
		s.Theme.Text(" 3:00 AM").FlexGrow("1"),
		s.Theme.Text(" 4:00 AM").FlexGrow("1"),
		s.Theme.Text(" 5:00 AM").FlexGrow("1"),
		s.Theme.Text(" 6:00 AM").FlexGrow("1"),
		s.Theme.Text(" 7:00 AM").FlexGrow("1"),
		s.Theme.Text(" 8:00 AM").FlexGrow("1"),
		s.Theme.Text(" 9:00 AM").FlexGrow("1"),
		s.Theme.Text("10:00 AM").FlexGrow("1"),
		s.Theme.Text("11:00 AM").FlexGrow("1"),
		s.Theme.Text("12:00 PM").FlexGrow("1"),
		s.Theme.Text(" 1:00 PM").FlexGrow("1"),
		s.Theme.Text(" 2:00 PM").FlexGrow("1"),
		s.Theme.Text(" 3:00 PM").FlexGrow("1"),
		s.Theme.Text(" 4:00 PM").FlexGrow("1"),
		s.Theme.Text(" 5:00 PM").FlexGrow("1"),
		s.Theme.Text(" 6:00 PM").FlexGrow("1"),
		s.Theme.Text(" 7:00 PM").FlexGrow("1"),
		s.Theme.Text(" 8:00 PM").FlexGrow("1"),
		s.Theme.Text(" 9:00 PM").FlexGrow("1"),
		s.Theme.Text("10:00 PM").FlexGrow("1"),
		s.Theme.Text("11:00 PM").FlexGrow("1").Color("lightgray"),
	).Color("lightgray").FontSize(browser.Size{Value: 10, Unit: browser.UnitPX}))

	for j, _ := range days {
		var rows []*browser.Node

		/*
			rows = append(rows, ui.HStack(
				s.Theme.Button(day.Format("Mon 2")).
					OnClick(func(d time.Time) func(browser.Dispatcher(EventDayClick{day})) {
						return func(_ dom.Event) {
							*s.Time = d
							*s.SelectedKey = "day"
						}
					}(day)),
			).JustifyContentCenter(),
			)
		*/

		rows = append(rows, gridSquare(s, -1, true, j == 0, false, false))

		for i := 0; i < 24; i++ {
			rows = append(rows, gridSquare(s, i, false, j == 0, i == 23, j == 6))
		}

		cols = append(cols, ui.VStack(rows...).FlexGrow("1"))
	}

	//	log.Printf("%s is a %s; month has %d days", day, day.Weekday(), len(days))

	//	return ui.VStack(s.Theme.Text("NONSENSE"))

	return ui.HStack(cols...).FlexGrow("1")
}

func gridSquare(s *State, hour int, rowStart, colStart, rowEnd, colEnd bool) *browser.Node {
	return ui.HStack(
		s.Theme.Textf("Hour %d", hour),
	).
		PaddingPX(10).
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

func border(th *ui.Theme) browser.Border {
	return browser.Border{
		Color: "lightgray", // TODO: use theme?
		Width: browser.Size{Value: 1, Unit: browser.UnitPX},
		Type:  browser.BorderSolid,
	}
}

func lineView(s *State, e *cal.EventItem) *browser.Node {
	return ui.VStack(
		s.Theme.Text(e.Name).OverflowHidden().
			FontSizeEM(1).
			//			MaxHeight(&browser.Size{Value: 1.2, Unit: browser.UnitEM}).
			Color("black"),
		ui.OnlyIf(e.HourSpecified,
			func() *browser.Node {
				return s.Theme.Text(e.Time.Format("3:04 PM")).Color("black")
			}),
	).
		PaddingPX(2).
		Background("lightgray").
		MarginTopPX(2).
		BorderRadiusPX(3).
		OnClick(browser.Dispatcher(EventEventClick{e})).
		Pointer()

}
