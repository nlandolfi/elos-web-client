package app

import (
	"log"
	"time"

	"github.com/nlandolfi/elos/web-client/components/calendar"
	"github.com/nlandolfi/elos/web-client/components/calendar/manager"
	"github.com/nlandolfi/elos/web-client/components/editor"
	"github.com/nlandolfi/elos/web-client/components/notes"
	"github.com/nlandolfi/spin/infra/ctzn"
	"github.com/nlandolfi/spin/infra/fs"
	"github.com/nlandolfi/spin/infra/key"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
	"github.com/spinsrv/web-client/components/icons"
	"github.com/spinsrv/web-client/components/login"
	"github.com/spinsrv/web-client/components/profile"
	"github.com/spinsrv/web-client/components/sidebar"
)

type State struct {
	Theme ui.Theme

	*key.PrivateKey
	LoginError string
	key.KeyServerHTTPClient

	SidebarState  sidebar.State
	SidebarHidden bool

	LoginState    login.State
	CalendarState calendar.State
	EditorState   editor.State
	ProfileState  profile.State
	NotesState    notes.State

	ClientVersion string
	LastWrittenAt time.Time
}

type EventInitialize struct{}
type EventLoginSuccess struct{}
type EventToggleSidebar struct{}

func (s *State) Handle(e browser.Event) {
	switch v := e.(type) {
	case EventInitialize:
		return
	case EventLoginSuccess:
		s.LoginState.Username = ""
		s.LoginState.Password = ""
		go browser.Dispatch(sidebar.EventItemClick{Target: &s.SidebarState, Item: *items[0]})
	case EventToggleSidebar:
		s.SidebarHidden = !s.SidebarHidden
	case login.EventLoginButtonClicked:
		go s.Login(s.LoginState.Username, s.LoginState.Password)
	case profile.EventSelectTheme:
		s.SetTheme(v.Key)
	case sidebar.EventItemClick:
		switch v.Key {
		case "logout":
			s.Logout()
		case "calendar":
			go s.CalendarState.Reload()
			s.SidebarState.Handle(e)
		case "editor":
			go s.EditorState.File.Reload()
			s.SidebarState.Handle(e)
		default:
			s.SidebarState.Handle(e)
		}
	case sidebar.EventItemHoverStart, sidebar.EventItemHoverEnd:
		s.SidebarState.Handle(e)
	case manager.EventOpenInEditor:
		s.OpenInEditor(
			ctzn.Name(s.CalendarState.ManagerState.Calendar.Citizen),
			fs.Path(s.CalendarState.ManagerState.Calendar.Path),
		)
	default:
	}

	s.LoginState.Handle(e)
	s.CalendarState.Handle(e)
	s.EditorState.Handle(e)
	s.NotesState.Handle(e)
}

func (s *State) Login(pu, pr string) {
	s.LoginError = "authenticating..."
	go browser.Dispatch(nil)
	defer func() { go browser.Dispatch(nil) }()

	resp := s.KeyServerHTTPClient.Temp(&key.TempRequest{
		Public:   pu,
		Private:  pr,
		Duration: 24 * time.Hour,
	})

	if resp.Error != "" {
		s.LoginError = resp.Error
		return
	}

	if resp.Key == nil {
		s.LoginError = "authentication failed"
		return
	}

	if s.PrivateKey == nil {
		s.PrivateKey = new(key.PrivateKey)
	}

	go browser.Dispatch(EventLoginSuccess{})

	s.PrivateKey.Key = *resp.Key
	s.PrivateKey.Private = resp.Private
	s.LoginError = ""
}

func (s *State) Rewire() {
	s.LoginState.Theme = &s.Theme
	s.LoginState.Status = &s.LoginError
	s.CalendarState.Rewire(&s.Theme, &s.PrivateKey)
	s.EditorState.Rewire(&s.Theme, &s.PrivateKey)
	s.SidebarState.Theme = &s.Theme
	s.ProfileState.Theme = &s.Theme
	s.ProfileState.PrivateKey = &s.PrivateKey
	s.NotesState.Rewire(&s.Theme, &s.PrivateKey)
}

func (s *State) Logout() {
	// *s.PrivateKey = nil
	// the above was previous memo, now just wipe the state to ensure no leaking data
	t := s.Theme
	v := s.ClientVersion
	*s = State{} // wipe state
	s.Theme = t
	s.ClientVersion = v
	s.Rewire()
}

func (s *State) OpenInEditor(c ctzn.Name, p fs.Path) {
	go s.EditorState.Open(c, p)
	s.SidebarState.SelectedKey = "editor"
	s.SidebarState.SelectedDisplay = "Editor"
}

func (s *State) SetTheme(to string) {
	switch to {
	case "light":
		s.Theme = ui.LightTheme
	case "dark":
		s.Theme = ui.DarkTheme
	default:
		s.Theme = ui.DefaultTheme
	}
}

func View(s *State) *browser.Node {
	return view(s).Background(s.Theme.BackgroundColor).FontFamily(s.Theme.FontFamily)
}

// todo: maybe don't re-declare
var items = []*sidebar.Item{
	&sidebar.Item{
		Key:     "calendar",
		Display: "Calendar",
	},
	&sidebar.Item{
		Key:     "notes",
		Display: "Notes",
	},
	&sidebar.Item{
		Key:     "editor",
		Display: "Editor",
	},
	&sidebar.Item{
		Key:     "profile",
		Display: "Profile",
	},
	&sidebar.Item{
		Key:     "logout",
		Display: "Logout",
	},
}

func view(s *State) *browser.Node {
	if s.PrivateKey == nil {
		return login.View(&s.LoginState).Background("black")
	}

	var view *browser.Node
	switch s.SidebarState.SelectedKey {
	case "calendar", "":
		view = calendar.View(&s.CalendarState)
	case "editor":
		view = editor.View(&s.EditorState)
	case "notes":
		view = notes.View(&s.NotesState)
	case "profile":
		view = profile.View(&s.ProfileState)
	default:
		log.Fatalf("unknown selected app: %q", s.SidebarState.SelectedKey)
	}

	return ui.VStack(
		header(s),
		ui.HStack(
			ui.OnlyIf(!s.SidebarHidden,
				func() *browser.Node {
					return sidebar.View(&s.SidebarState, items).Width(
						browser.Size{Value: 100, Unit: browser.UnitPX},
					).PaddingPX(10).BorderRight(border)
				},
			),
			view.FlexGrow("1"),
		).FlexGrow("1"),
	).Height(
		browser.Size{Value: 100, Unit: browser.UnitVH},
	).OverflowScroll()

}

func header(s *State) *browser.Node {
	return ui.HStack(
		ui.HStack(
			icons.Trademark(s.Theme.BackgroundColor, 30),
			ui.OnlyIf(s.SidebarHidden,
				func() *browser.Node { return s.Theme.Text(s.SidebarState.SelectedDisplay) },
			),
		).OnClickDispatch(EventToggleSidebar{}).
			FlexGrow("1").AlignItemsCenter(),
		s.Theme.Textf("v%s", s.ClientVersion).
			FontSize(browser.Size{Value: 10, Unit: browser.UnitPX}).
			MarginRight(browser.Size{Value: 15, Unit: browser.UnitPX}),
	).BorderBottom(browser.Border{
		Width: browser.Size{Value: 1, Unit: browser.UnitPX},
		Type:  browser.BorderSolid,
		Color: "lightgray",
	},
	).AlignItemsCenter().CursorPointer()
}

var border = browser.Border{
	Width: browser.Size{Value: 1, Unit: browser.UnitPX},
	Type:  browser.BorderSolid,
	Color: "lightgray",
}
