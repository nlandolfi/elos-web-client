package manager

import (
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"github.com/spinsrv/browser/ui"
)

type CalendarReference struct {
	Citizen string
	Path    string
}

type EventOpenInEditor struct{}
type EventReloadCalendar struct{}

type State struct {
	Theme *ui.Theme

	Calendar CalendarReference
}

func View(s *State) *browser.Node {
	return ui.VStack(
		s.Theme.Text("The loaded calendar:"),
		ui.HStack(
			s.Theme.TextInput(&(s.Calendar.Citizen)).
				Placeholder("citizen").
				OnKeyDown(func(e dom.Event) {
					if e.KeyCode() == 13 { // enter
						go browser.Dispatch(EventReloadCalendar{}) // todo change name of action
					}
				}),
			s.Theme.TextInput(&(s.Calendar.Path)).
				Placeholder("path").
				OnKeyDown(func(e dom.Event) {
					if e.KeyCode() == 13 { // enter
						go browser.Dispatch(EventReloadCalendar{}) // todo change name of action
					}
				}),
			s.Theme.Button("Open in editor").OnClickDispatch(EventOpenInEditor{}),
		).FlexGrow("1"),
	)
}
