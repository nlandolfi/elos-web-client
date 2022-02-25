package manager

import (
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type NotesRoot struct {
	Citizen string
	Path    string
}

type State struct {
	Theme *ui.Theme

	NotesRoot
}

func View(s *State) *browser.Node {
	return ui.VStack(
		s.Theme.Text("The notes root:"),
		ui.HStack(
			s.Theme.TextInput(&(s.NotesRoot.Citizen)).
				Placeholder("citizen"),
			s.Theme.TextInput(&(s.NotesRoot.Path)).
				Placeholder("path"),
		).FlexGrow("1"),
	)
}
