package editor

import (
	"time"

	"github.com/nlandolfi/spin/apps/cal"
	"github.com/nlandolfi/spin/infra/key"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme       *ui.Theme        `json:"-"`
	PrivateKey  **key.PrivateKey `json:"-"`
	Time        *time.Time       `json:"-"`
	SelectedKey *string          `json:"-"`

	Status    string
	debugging bool
	Event     *cal.EventItem
}

type EventCreateEvent struct{}
type EventAddExclude struct{}
type EventDropExcludes struct{ index int }
type EventToggleDebugging struct{}
type EventValidate struct{}
type EventSave struct{}
type EventCancel struct{}
type EventDelete struct{ ID string }
type EventBack struct{}
type EventSaved struct{}

func (s *State) Handle(e browser.Event) {
	switch e := e.(type) {
	case EventAddExclude:
		s.Event.Excludes = append(s.Event.Excludes, time.Now())
	case EventToggleDebugging:
		s.debugging = !s.debugging
	case EventCancel:
		s.Event = nil
		s.Status = ""
	case EventDropExcludes:
		var es []time.Time
		for j, o := range s.Event.Excludes {
			if j == e.index {
				continue
			}

			es = append(es, o)
		}
		s.Event.Excludes = es
	case EventValidate:
		s.Status = "validating..."
		if s.Event.Time.IsZero() {
			s.Status = "need a nonzero time"
			return
		}
		go browser.Dispatch(EventSave{})
	case EventSaved:
		*s.SelectedKey = "day"
		*s.Time = s.Event.Time
		s.Event = nil
		s.Status = ""
	}
}

func (s *State) Rewire(th *ui.Theme, k **key.PrivateKey) {
	s.Theme = th
	s.PrivateKey = k
}

func View(s *State) *browser.Node {
	if s.Event == nil {
		return s.Theme.Button("Create New Event").OnClickDispatch(EventCreateEvent{})
	}

	return s.Theme.Card(ui.VStack(
		ui.HStack(
			// Header
			ui.If(s.Event.ID == "",
				func() *browser.Node { return s.Theme.Text("Creating event:") },
				func() *browser.Node { return s.Theme.Textf("Editing event %q:", s.Event.ID) },
			).MarginRightPX(10),
			s.Theme.Text(s.Status),
		),

		s.Theme.TextInput(&(s.Event.Name)).Placeholder("Name"),
		s.Theme.TextArea(&(s.Event.Details)).Placeholder("Details"),
		ui.VStack(
			s.Theme.DateInput(&(s.Event.Time)),

			ui.HStack(
				s.Theme.Text("Hours?"),
				s.Theme.Toggle(&s.Event.HourSpecified),
			),

			ui.OnlyIf(s.Event.HourSpecified,
				func() *browser.Node { return s.Theme.TimeInput(&s.Event.Time) },
			),

			ui.VStack(
				ui.HStack(
					s.Theme.Text("Recurs?"),
					s.Theme.Toggle(&s.Event.Recurs),
				),

				ui.OnlyIf(s.Event.Recurs,
					func() *browser.Node {
						return ui.VStack(
							ui.HStack(
								s.Theme.Textf("This event recurs"),
								s.Theme.TextInput(&(s.Event.Frequency)),
								s.Theme.Textf("at an interval of %d", s.Event.Interval),
							),

							ui.HStack(
								s.Theme.Text("Interval?"),
								s.Theme.Toggle(&s.Event.IntervalSpecified),
							),
							ui.OnlyIf(s.Event.IntervalSpecified,
								func() *browser.Node { return s.Theme.NumberInput(&s.Event.Interval).Min(1).Max(12) },
							),
							ui.HStack(
								s.Theme.Text("Until?"),
								s.Theme.Toggle(&s.Event.UntilSpecified),
							),
							ui.OnlyIf(s.Event.UntilSpecified,
								func() *browser.Node { return s.Theme.DateInput(&(s.Event.Until)) },
							),
							ui.HStack(
								s.Theme.Text("Exclude:"),
								s.Theme.Button("Add Date Exception").OnClickDispatch(EventAddExclude{}),
							),
							ui.OnlyIf(len(s.Event.Excludes) > 0,
								func() *browser.Node { return excludes(s.Theme, s.Event) },
							),
						)
					},
				),
			),

			ui.HStack(
				ui.VStack(
					s.Theme.Button("Debug").OnClickDispatch(EventToggleDebugging{}),
				),
				s.Theme.Button("Cancel").OnClickDispatch(EventCancel{}),
				ui.OnlyIf(s.Event.ID != "",
					func() *browser.Node {
						return s.Theme.Button("Delete").OnClickDispatch(EventDelete{s.Event.ID})
					}),
				s.Theme.Button("Save").OnClickDispatch(EventValidate{}),
				s.Theme.Button("Back to previous").OnClickDispatch(EventBack{}),
			),

			ui.OnlyIf(s.debugging,
				func() *browser.Node { return s.Theme.Textf("%+v", s.Event) },
			),
		),
	))
}

func excludes(t *ui.Theme, e *cal.EventItem) *browser.Node {
	var views []*browser.Node

	views = append(views,
		t.Text("Exclusions applied:"),
	)

	for i := range e.Excludes {
		views = append(views, ui.HStack(
			t.DateInput(&e.Excludes[i]),
			t.Button("X").OnClickDispatch(EventDropExcludes{i}),
		))
	}

	return ui.VStack(
		views...,
	)
}
