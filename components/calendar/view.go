package calendar

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/nlandolfi/elos/web-client/components/calendar/day"
	"github.com/nlandolfi/elos/web-client/components/calendar/editor"
	"github.com/nlandolfi/elos/web-client/components/calendar/inspector"
	"github.com/nlandolfi/elos/web-client/components/calendar/manager"
	"github.com/nlandolfi/elos/web-client/components/calendar/month"
	"github.com/nlandolfi/elos/web-client/components/calendar/week"
	"github.com/nlandolfi/elos/web-client/components/calendar/year"
	"github.com/nlandolfi/elos/web-client/components/selector"
	"github.com/nlandolfi/spin/apps/cal"
	"github.com/nlandolfi/spin/infra/key"
	feditor "github.com/nlandolfi/spin/web/elos/components/editor"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme      *ui.Theme        `json:"-"`
	PrivateKey **key.PrivateKey `json:"-"`

	Time time.Time

	SelectorState selector.State

	DayState     day.State
	MonthState   month.State
	WeekState    week.State
	YearState    year.State
	EditorState  editor.State
	ManagerState manager.State
	// TableState       table.State

	// Inspector
	InspectorState   inspector.State
	InspectorVisible bool
	InspectedEvent   *cal.EventItem

	LastSelectedItem selector.Item

	CalendarFile feditor.File
	EventItems   []*cal.EventItem
}

type EventEditEvent struct{ *cal.EventItem }
type EventReloadEvents struct{}

func (s *State) Handle(e browser.Event) {
	switch e := e.(type) {
	case week.EventEventClick:
		s.InspectorVisible = true
		s.InspectedEvent = e.EventItem
	case inspector.EventEditEvent:
		go browser.Dispatch(EventEditEvent{s.InspectedEvent})
	case manager.EventReloadCalendar, EventReloadEvents:
		go s.reloadEvents()
	case EventEditEvent:
		s.editEvent(e.EventItem)
		go browser.Dispatch(selector.EventItemClick{
			Target: &s.SelectorState,
			Item:   *SelectorItems[4],
		})
	case editor.EventCreateEvent:
		i := new(cal.EventItem)
		s.EventItems = append(s.EventItems, i)
		s.editEvent(i)
	case editor.EventCancel:
		go browser.Dispatch(selector.EventItemClick{
			Target: &s.SelectorState,
			Item:   s.LastSelectedItem,
		})
	case editor.EventSave:
		go s.save()
	case editor.EventBack:
		go browser.Dispatch(selector.EventItemClick{
			Target: &s.SelectorState,
			Item:   s.LastSelectedItem,
		})
	case editor.EventDelete:
		var es []*cal.EventItem
		for _, i := range s.EventItems {
			if e.ID != i.ID {
				es = append(es, i)
			}
		}
		s.EventItems = es
		go s.save()
	}
	s.SelectorState.Handle(e)
	s.DayState.Handle(e)
	s.MonthState.Handle(e)
	s.WeekState.Handle(e)
	s.YearState.Handle(e)
	s.InspectorState.Handle(e)
	s.EditorState.Handle(e)
}

func (s *State) Rewire(th *ui.Theme, k **key.PrivateKey) {
	s.SetTheme(th)
	s.SetPrivateKey(k)

	if s.SelectorState.SelectedKey == "" {
		s.SelectorState.SelectedKey = "month"
		s.SelectorState.SelectedDisplay = "Month"
	}

	if s.Time.IsZero() {
		s.Time = time.Now()
	}

	//	s.TableState.EventItems = s.EventItems

	s.DayState.Time = &s.Time
	s.DayState.EventItems = s.EventItems
	s.DayState.DispatchEditEvent = func(e *cal.EventItem) {
		go browser.Dispatch(EventEditEvent{e})
	}

	s.WeekState.Time = &s.Time
	s.WeekState.EventItems = s.EventItems
	s.WeekState.SelectedKey = &s.SelectorState.SelectedKey

	s.MonthState.Time = &s.Time
	s.MonthState.EventItems = s.EventItems
	s.MonthState.SelectedKey = &s.SelectorState.SelectedKey
	s.MonthState.InspectEvent = s.inspectEvent

	s.YearState.Time = &s.Time
	s.YearState.EventItems = s.EventItems
	s.YearState.SelectedKey = &s.SelectorState.SelectedKey

	s.InspectorState.Visible = &s.InspectorVisible
	s.InspectorState.EventItem = &s.InspectedEvent

	s.EditorState.Rewire(th, k)
	s.EditorState.Time = &s.Time
	s.EditorState.SelectedKey = &s.SelectorState.SelectedKey
}

func (s *State) SetTheme(th *ui.Theme) {
	s.Theme = th
	s.SelectorState.Theme = th
	//	s.TableState.SetTheme(th)
	s.MonthState.SetTheme(th)
	s.WeekState.SetTheme(th)
	s.DayState.SetTheme(th)
	s.YearState.SetTheme(th)
	s.InspectorState.SetTheme(th)
	s.ManagerState.Theme = th
}

func (s *State) SetPrivateKey(k **key.PrivateKey) {
	s.PrivateKey = k
}

var SelectorItems = []*selector.Item{
	&selector.Item{
		Key:     "day",
		Display: "Day",
	},
	&selector.Item{
		Key:     "week",
		Display: "Week",
	},
	&selector.Item{
		Key:     "month",
		Display: "Month",
	},
	&selector.Item{
		Key:     "year",
		Display: "Year",
	},
	//	&selector.Item{
	//		Key:     "table",
	//		Display: "Table",
	//	},
	&selector.Item{
		Key:     "editor",
		Display: "Editor",
	},
	&selector.Item{
		Key:     "manager",
		Display: "Manager",
	},
}

func View(s *State) *browser.Node {
	return ui.ZStack(
		ui.VStack(
			ui.HStack(
				selector.View(&s.SelectorState, SelectorItems),
				s.Theme.Button("Reload").OnClickDispatch(EventReloadEvents{}).PaddingPX(9).PaddingBottomPX(3),
				s.Theme.Text(s.CalendarFile.Status),
			).AlignItemsCenter(),
			view(s).PaddingPX(10),
		),
		ui.OnlyIf(s.InspectorVisible && (s.SelectorState.SelectedKey == "week" || s.SelectorState.SelectedKey == "month"),
			func() *browser.Node { return inspector.View(&s.InspectorState) },
		).PositionAbsolute(), // for positioning
	).PositionRelative() // relative for the inspector
}

func view(s *State) *browser.Node {
	switch s.SelectorState.SelectedKey {
	//	case "table":
	//		return table.View(&s.TableState)
	case "day", "":
		return day.View(&s.DayState)
	case "week":
		return week.View(&s.WeekState)
	case "month":
		return month.View(&s.MonthState)
	case "year":
		return year.View(&s.YearState)
	case "editor":
		return editor.View(&s.EditorState)
	case "manager":
		return manager.View(&s.ManagerState)
	default:
		panic(fmt.Sprintf("unknown selected state: %v", s.SelectorState.SelectedKey))
	}
}

func (s *State) inspectEvent(e *cal.EventItem) {
	s.InspectorVisible = true
	s.InspectedEvent = e
	s.InspectorState.EventItem = &s.InspectedEvent
}

func (s *State) save() {
	log.Print("writing events")
	var b bytes.Buffer

	if err := cal.WriteEvents(&b, s.EventItems); err != nil {
		s.CalendarFile.Status = err.Error() // TODO: don't use this status field
	}
	log.Print("and saving")

	s.CalendarFile.Text = b.String()
	s.CalendarFile.Save()
	s.reloadEventItems()
	go browser.Dispatch(editor.EventSaved{})
}

func (s *State) Reload() {
	s.reloadEvents()
}

func (s *State) reloadEventItems() {
	b := bytes.NewBufferString(s.CalendarFile.Text)

	es, err := cal.ParseEvents(b)
	if err != nil {
		s.CalendarFile.Status = err.Error()
		return
	}
	for i, e := range es {
		e.ID = fmt.Sprintf("%d-%d", s.CalendarFile.ShadowSequence, i)
	}
	s.EventItems = es
	s.Rewire(s.Theme, s.PrivateKey)
}

func (s *State) reloadEvents() {
	s.InspectorVisible = false
	s.InspectedEvent = nil
	s.EditorState.Event = nil

	s.CalendarFile.PrivateKey = s.PrivateKey
	s.CalendarFile.Citizen = s.ManagerState.Calendar.Citizen
	s.CalendarFile.Path = s.ManagerState.Calendar.Path
	s.CalendarFile.Reload()

	s.reloadEventItems()

	/*
		s.Error = "loading..."
		s.CalendarFile.Reload()
		go browser.Dispatch(nil)
		defer func() { go browser.Dispatch(nil) }()

		s.CalendarFile.PrivateKey = s.PrivateKey
		s.CalendarFile.Citizen = s.ManagerState.Calendar.Citizen

		k := *s.PrivateKey
		s.Error = "loading"
		go browser.Dispatch(nil)
		defer func() { go browser.Dispatch(nil) }()
		c := new(cal.CalServerHTTPClient)
		p := fs.Path(s.ManagerState.Calendar.Path)
		if s.Path != "" {
			p = fs.Path(s.Path)
		}

		resp := c.All(&cal.CalAllRequest{
			Public: string(k.Name), Private: k.Private,
			Citizen: ctzn.Name(s.ManagerState.Calendar.Citizen), Path: p,
		})

		if resp.Error != "" {
			s.Error = resp.Error
			return
		}

		sort.Sort(cal.ByTime{resp.Events})

		s.EventItems = resp.Events
		s.Error = ""
		// hack to rewire events
		s.Rewire(s.Theme, s.PrivateKey) // TODO
	*/
}

func (s *State) editEvent(e *cal.EventItem) {
	s.LastSelectedItem = selector.Item{
		Key:     s.SelectorState.SelectedKey,
		Display: s.SelectorState.SelectedDisplay,
	}
	s.EditorState.Event = e
}
