package year

import (
	"log"
	"time"

	"github.com/nlandolfi/elos/web-client/components/calendar/ltr"
	"github.com/nlandolfi/elos/web-client/components/calendar/mgrid"
	"github.com/nlandolfi/spin/apps/cal"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme       *ui.Theme        `json:"-"`
	Time        *time.Time       `json:"-"`
	EventItems  []*cal.EventItem `json:"-"` // this is a pointer
	SelectedKey *string          `json:"-"`

	hoveredMonth time.Time
	hoveredDay   time.Time
}

func (s *State) Handle(e browser.Event) {
	switch e := e.(type) {
	case EventMonthClick:
		*s.SelectedKey = "month"
		*s.Time = e.Time
	case EventMonthHoverStart:
		s.hoveredMonth = e.Time
	case EventMonthHoverLeave:
		if s.hoveredMonth.Equal(e.Time) {
			s.hoveredMonth = time.Time{}
		}
	case EventDayClick:
		log.Print("click day")
		*s.SelectedKey = "day"
		*s.Time = e.Time
	case EventDayHoverStart:
		s.hoveredDay = e.Time
	case EventDayHoverLeave:
		if s.hoveredDay.Equal(e.Time) {
			s.hoveredDay = time.Time{}
		}
	case EventDecrementYear:
		*s.Time = time.Date(
			s.Time.Year()-1, s.Time.Month(), s.Time.Day(),
			s.Time.Hour(), s.Time.Minute(), s.Time.Second(), s.Time.Nanosecond(),
			s.Time.Location(),
		)
	case EventSetTime:
		*s.Time = e.Time
	case EventIncrementYear:
		*s.Time = time.Date(
			s.Time.Year()+1, s.Time.Month(), s.Time.Day(),
			s.Time.Hour(), s.Time.Minute(), s.Time.Second(), s.Time.Nanosecond(),
			s.Time.Location(),
		)
	}
}

type EventMonthClick struct{ time.Time }
type EventMonthHoverStart struct{ time.Time }
type EventMonthHoverLeave struct{ time.Time }
type EventDayClick struct{ time.Time }
type EventDayHoverStart struct{ time.Time }
type EventDayHoverLeave struct{ time.Time }
type EventDecrementYear struct{}
type EventSetTime struct{ time.Time }
type EventIncrementYear struct{}

func (s *State) SetTheme(t *ui.Theme) {
	s.Theme = t
}

func View(s *State) *browser.Node {
	m, n := 3, 4 // 3 * 4 = 12 months

	var rows []*browser.Node = make([]*browser.Node, m)

	for i := 0; i < m; i++ {
		var cols []*browser.Node = make([]*browser.Node, n)

		for j := 0; j < n; j++ {
			sentinel := time.Date(s.Time.Year(), time.Month((i*n)+j+1), 1, 0, 0, 0, 0, s.Time.Location())
			lcache := sentinel.String() // lcache for listener cache

			cols[j] = ui.VStack(
				s.Theme.Textf("%s", sentinel.Month()).
					FontSizeEM(1.2).
					TextAlignCenter().
					//TextAlignRight().
					//PaddingRightPX(20).
					Pointer().
					OnClickCached(lcache, browser.Dispatcher(EventMonthClick{sentinel})).
					OnlyIf(cal.SameDay(s.hoveredMonth, sentinel),
						func(n *browser.Node) *browser.Node {
							return n.Background("lightgray").BorderRadiusPX(5) // TODO use theme
						}).
					OnMouseEnter(browser.Dispatcher(EventMonthHoverStart{sentinel})).
					OnMouseLeave(browser.Dispatcher(EventMonthHoverLeave{sentinel})), // TODO: cache these again
				monthGrid(s, sentinel),
			).PaddingPX(20).FlexGrow("1")
		}

		rows[i] = ui.HStack(cols...).FlexGrow("1")
	}

	return ui.VStack(
		header(s),
		ui.VStack(rows...).FlexGrow("1"),
	).FlexGrow("1")
}

func header(s *State) *browser.Node {
	return ui.HStack(
		s.Theme.Text(s.Time.Format("2006")).
			FontSizeEM(1.5).
			MarginLeftPX(10),
		/*
			ui.OnlyIf(s.Time.Year() == time.Now().Year(),
				func() *browser.Node {
					return s.Theme.Text(s.Time.Format("(this year)")).
						FontSizeEM(1).
						MarginLeftPX(10)
				}),
		*/
		ui.Spacer(),
		ltr.View(&ltr.State{
			Theme:        s.Theme,
			OnClickPrev:  browser.Dispatcher(EventDecrementYear{}),
			OnClickToday: browser.Dispatcher(EventSetTime{time.Now()}),
			OnClickNext:  browser.Dispatcher(EventIncrementYear{}),
		}),
	).AlignItemsCenter()
}

func monthGrid(s *State, sentinel time.Time) *browser.Node {
	return mgrid.Grid(&mgrid.GridSpec{
		PadTo6Weeks: true,
		Time:        sentinel,
		Head: func(days []time.Time, i int) *browser.Node {
			return ui.VStack(
				ui.VSpace(browser.Size{Value: 3, Unit: browser.UnitPX}),
				s.Theme.Text(days[i].Weekday().String()[:1]).
					PaddingPX(3).
					Color("lightgray").
					TextAlignCenter(),
				ui.VSpace(browser.Size{Value: 3, Unit: browser.UnitPX}),
			).
				FlexGrow("1").
				FlexBasis("0px")
		},
		Cell: func(days []time.Time, i, j int, t time.Time) *browser.Node {
			// this is the value which, so long as it doesn't change, none of the
			// handlers needs to change - NCL 2/2/2022
			lKeyCache := t.String()

			number := s.Theme.Textf("%d", t.Day()).
				OnlyIf(t.Month() != sentinel.Month(),
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
				PaddingPX(3).
				FlexGrow("1").
				FlexBasis("0px").
				Pointer().
				TextAlignCenter().
				OnClickDispatchCached(lKeyCache, EventDayClick{t}).
				OnMouseEnter(browser.Dispatcher(EventDayHoverStart{t})).
				OnMouseLeave(browser.Dispatcher(EventDayHoverLeave{t}))

			var dots []*browser.Node
			for _, e := range s.EventItems {
				if e.ShouldDisplayOnDay(t) {
					dots = append(dots, s.Theme.Text("*"))
				}
			}

			return ui.VStack(
				number.OnlyIf(len(dots) > 0, func(n *browser.Node) *browser.Node {
					return n.FontWeight("700") //TextDecorationUnderline()
				}),
			).FlexGrow("1").
				FlexBasis("0px")
		},
	}).FlexGrow("1")
}
