package inspector

import (
	"log"

	"github.com/nlandolfi/spin/apps/cal"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme *ui.Theme

	Visible   *bool
	EventItem **cal.EventItem

	CurrentX, CurrentY int
	dragging           bool
}

func (s *State) Handle(e browser.Event) {
	switch e := e.(type) {
	case EventDragBegin:
		log.Print("drag begin")
		s.dragging = true
	case EventDragHault:
		s.dragging = false
	case EventDragMove:
		if s.dragging {
			s.CurrentX += e.DiffX
			s.CurrentY += e.DiffY
		}
	case EventClose:
		*s.Visible = false
	}
}

type EventDragBegin struct{}
type EventDragHault struct{}
type EventDragMove struct{ DiffX, DiffY int }
type EventClose struct{}
type EventEditEvent struct{}

func (s *State) SetTheme(t *ui.Theme) {
	s.Theme = t
}

func View(s *State) *browser.Node {
	return s.Theme.Card(
		// this is the header
		ui.HStack(
			/*
				ui.VStack(
					s.Theme.Text(LeftChevron).FontSizeEM(2),
				).FlexGrow("1").JustifyContentCenter(),
			*/
			cardBodyView(s).FlexGrow("1"),
			/*
				ui.VStack(
					s.Theme.Text(RightChevron),
				).FlexGrow("1").JustifyContentCenter(),
			*/
		).PaddingLeftPX(30).
			PaddingRightPX(30).
			PaddingBottomPX(30).
			MaxWidthPX(300),
	).MaxWidth(browser.Size{Value: 500, Unit: browser.UnitPX}).
		PositionAbsolute().
		LeftPX(float64(s.CurrentX)).
		TopPX(float64(s.CurrentY))
	//	.PaddingPX(0).MarginPX(0)
}

func cardBodyView(s *State) *browser.Node {
	if s.EventItem == nil {
		return s.Theme.Text("nil")
	}
	item := *s.EventItem
	if item == nil {
		return s.Theme.Text("nil")
	}
	return ui.VStack(
		ui.HStack(
			ui.Spacer(),
			s.Theme.Text("Event Inspector").
				MarginPX(10).
				FlexGrow("1").
				CursorMove().
				OnMouseDown(browser.Dispatcher(EventDragBegin{})).
				OnMouseUp(browser.Dispatcher(EventDragHault{})).
				OnMouseMove(func(e dom.Event) {
					e.PreventDefault()
					e.StopPropagation()

					go browser.Dispatch(EventDragMove{DiffX: e.MovementX(), DiffY: e.MovementY()})
				}),
			s.Theme.Button("x").OnClickDispatch(EventClose{}),
		),
		ui.HStack(
			s.Theme.Text(item.Name).
				FontSizeEM(1.2),
		),
		s.Theme.Text(item.Time.Format("2 Jan 2006")).FontSizeEM(1.0),
		ui.OnlyIf(item.HourSpecified,
			func() *browser.Node { return s.Theme.Text(item.Time.Format("3:04 PM")) },
		),
		ui.If(len(item.Details) == 0,
			func() *browser.Node {
				return s.Theme.Text("No details...").Color("lightgray").FontSizeEM(0.5) // TODO: pull color from theme
			},
			func() *browser.Node {
				return s.Theme.Text(item.Details)
			},
		),
		ui.If(item.Recurs,
			func() *browser.Node {
				return s.Theme.Text(item.RecurString())
			},
			func() *browser.Node {
				return s.Theme.Text("No recurrence...").Color("lightgray").FontSizeEM(0.5) // TODO: pull color from theme
			},
		),
		s.Theme.Button("Edit").OnClickDispatch(EventEditEvent{}),
	)
}

const LeftChevron = "‹"
const RightChevron = "›"

// https://www.w3schools.com/howto/howto_js_draggable.asp

/*
// Make the DIV element draggable:
dragElement(document.getElementById("mydiv"));

function dragElement(elmnt) {
  var pos1 = 0, pos2 = 0, pos3 = 0, pos4 = 0;
  if (document.getElementById(elmnt.id + "header")) {
    // if present, the header is where you move the DIV from:
    document.getElementById(elmnt.id + "header").onmousedown = dragMouseDown;
  } else {
    // otherwise, move the DIV from anywhere inside the DIV:
    elmnt.onmousedown = dragMouseDown;
  }

  function dragMouseDown(e) {
    e = e || window.event;
    e.preventDefault();
    // get the mouse cursor position at startup:
    pos3 = e.clientX;
    pos4 = e.clientY;
    document.onmouseup = closeDragElement;
    // call a function whenever the cursor moves:
    document.onmousemove = elementDrag;
  }

  function elementDrag(e) {
    e = e || window.event;
    e.preventDefault();
    // calculate the new cursor position:
    pos1 = pos3 - e.clientX;
    pos2 = pos4 - e.clientY;
    pos3 = e.clientX;
    pos4 = e.clientY;
    // set the element's new position:
    elmnt.style.top = (elmnt.offsetTop - pos2) + "px";
    elmnt.style.left = (elmnt.offsetLeft - pos1) + "px";
  }

  function closeDragElement() {
    // stop moving when mouse button is released:
    document.onmouseup = null;
    document.onmousemove = null;
  }
}
*/
