package selector

import (
	"fmt"

	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme  *ui.Theme `json:"-"`
	Prefix string

	hoveredKey      string
	SelectedKey     string
	SelectedDisplay string
}

func (s *State) Handle(e browser.Event) {
	switch e := e.(type) {
	case EventItemClick:
		if e.Target != s {
			return
		}

		s.SelectedKey = e.Key
		s.SelectedDisplay = e.Display
	case EventItemHoverStart:
		if e.Target != s {
			return
		}

		s.hoveredKey = e.Key
	case EventItemHoverEnd:
		if e.Target != s {
			return
		}

		if s.hoveredKey == e.Key {
			s.hoveredKey = ""
		}
	default:
	}
}

type Item struct {
	Key     string
	Display string
}

// structs for extensibility
type EventItemHoverStart struct {
	Target *State
	Item
}
type EventItemHoverEnd struct {
	Target *State
	Item
}
type EventItemClick struct {
	Target *State
	Item
}

func View(s *State, items []*Item) *browser.Node {
	views := make([]*browser.Node, len(items))

	for index, item := range items {
		views[index] = itemView(s, index, item)
	}

	return ui.HStack(views...)
}

func itemView(s *State, index int, item *Item) *browser.Node {
	return s.Theme.Text(item.Display).
		PaddingPX(9).
		PaddingBottomPX(3).
		CursorPointer().
		OnlyIf(s.hoveredKey == item.Key, func(n *browser.Node) *browser.Node {
			return n.BorderBottom(browser.Border{
				Color: s.Theme.HoverBackgroundColor,
				Type:  browser.BorderSolid,
				Width: browser.Size{1, browser.UnitPX},
			})
		}).
		OnlyIf(s.SelectedKey == item.Key || (s.SelectedKey == "" && index == 0), func(n *browser.Node) *browser.Node {
			return n.Color("blue").BorderBottom(browser.Border{
				Color: "blue",
				Type:  browser.BorderSolid,
				Width: browser.Size{1, browser.UnitPX},
			})
			//				FontWeight("500")
		}).
		OnClickDispatch(EventItemClick{Target: s, Item: *item}).
		OnMouseEnterDispatch(EventItemHoverStart{Target: s, Item: *item}).
		OnMouseLeaveDispatch(EventItemHoverEnd{Target: s, Item: *item}).
		ID(fmt.Sprintf("sidebar-%s", item.Key))
}
