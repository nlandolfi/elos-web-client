package ltr

import (
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme                                  *ui.Theme
	OnClickPrev, OnClickNext, OnClickToday dom.EventHandler
}

const LeftChevron = "‹"
const RightChevron = "›"

func maybe(f dom.EventHandler) dom.EventHandler {
	return func(e dom.Event) {
		if f != nil {
			f(e)
		}
	}
}

func View(s *State) *browser.Node {
	return ui.HStack(
		s.Theme.Button(LeftChevron).OnClick(maybe(s.OnClickPrev)),
		s.Theme.Button("Today").OnClick(maybe(s.OnClickToday)),
		s.Theme.Button(RightChevron).OnClick(maybe(s.OnClickNext)),
	).MarginLeftPX(20)
}
