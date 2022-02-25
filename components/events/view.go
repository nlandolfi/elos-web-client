package events

import (
	"log"

	"github.com/nlandolfi/elos/web-client/components/calendar/month"
	"github.com/nlandolfi/elos/web-client/components/selector"
	"github.com/nlandolfi/spin/apps/cal"
	"github.com/nlandolfi/spin/infra/fs"
	"github.com/nlandolfi/spin/infra/key"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme      *ui.Theme        `json:"-"`
	PrivateKey **key.PrivateKey `json:"-"`

	Path  string
	Error string

	SelectorState selector.State
	MonthState    month.State

	EventItems []*cal.EventItem
}

func (s *State) RegisterTheme(th *ui.Theme) {
	s.Theme = th
}

func View() *browser.Node {
	return oldview(s)
}

func oldview(s *State) *browser.Node {
	var views []*browser.Node

	for _, item := range s.EventItems {
		views = append(views, eventCard(s, item))
	}

	return ui.VStack(
		ui.HStack(
			s.Theme.TextInput(&s.Path).Placeholder("path to spin file...").FlexGrow("1"),
			s.Theme.Button("Add Event").OnClick(func(e dom.Event) {
				log.Print("Add Event")
			}),
			s.Theme.Button("Refresh").OnClick(func(e dom.Event) {
				go loadItems(s)
			}),
		),
		s.Theme.Text(s.Error),
		ui.VStack(views...).Padding(10),
	)
}

func eventCard(s *State, item *cal.EventItem) *browser.Node {
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

func loadItems(s *State) {
	k := *s.PrivateKey
	s.Error = "loading"
	c := new(cal.CalServerHTTPClient)

	resp := c.All(&cal.CalAllRequest{
		Public: string(k.Name), Private: k.Private,
		Citizen: k.Citizen, Path: fs.Path(s.Path),
	})

	if resp.Error != "" {
		s.Error = resp.Error
		return
	}

	s.EventItems = resp.Events
	s.Error = ""
}
