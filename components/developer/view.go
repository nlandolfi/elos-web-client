package developer

import (
	"github.com/nlandolfi/elos/web-client/components/app"
	"github.com/spinsrv/browser"
)

// UNUSED THIS FILE
type State struct {
	AppState  app.State
	Rendering string
}

func View(s *State) *browser.Node {
	return ui.VStack(
		app.View(&s.AppState),
		ui.TextInput(&s.Rendering),
	)
}
