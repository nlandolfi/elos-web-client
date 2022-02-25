package notes

import (
	"bytes"
	"fmt"

	"github.com/nlandolfi/elos/web-client/components/notes/canvas"
	"github.com/nlandolfi/elos/web-client/components/notes/manager"
	"github.com/nlandolfi/elos/web-client/components/notes/markdown"
	"github.com/nlandolfi/elos/web-client/components/notes/note"
	"github.com/nlandolfi/elos/web-client/components/notes/prototype"
	"github.com/nlandolfi/elos/web-client/components/selector"
	"github.com/nlandolfi/spin/infra/ctzn"
	"github.com/nlandolfi/spin/infra/fs"
	"github.com/nlandolfi/spin/infra/key"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme      *ui.Theme        `json:"-"`
	PrivateKey **key.PrivateKey `json:"-"`

	SelectorState selector.State

	DirEntries []*fs.DirEntry
	Status     string

	MarkdownState  markdown.State
	ManagerState   manager.State
	PrototypeState prototype.State
	CanvasState    canvas.State
}

func (s *State) Rewire(t *ui.Theme, k **key.PrivateKey) {
	s.Theme = t
	s.PrivateKey = k
	s.SelectorState.Theme = t
	s.ManagerState.Theme = t
	s.MarkdownState.Theme = t
	s.PrototypeState.Theme = t
	s.CanvasState.Theme = t

	if s.PrototypeState.Root == nil {
		n, err := note.Parse(bytes.NewBufferString(s.PrototypeState.Raw))
		if err != nil {
			s.PrototypeState.Root = note.Document("", "", note.DocumentArticle)

		} else {
			s.PrototypeState.Root = n
		}
	}
}

var items = []*selector.Item{
	&selector.Item{
		Key:     "main",
		Display: "Main",
	},
	&selector.Item{
		Key:     "markdown",
		Display: "Markdown",
	},
	&selector.Item{
		Key:     "prototype",
		Display: "Prototype",
	},
	&selector.Item{
		Key:     "canvas",
		Display: "Canvas",
	},
	&selector.Item{
		Key:     "manager",
		Display: "Manager",
	},
}

type EventReloadNotes struct{}

func (s *State) Handle(e browser.Event) {
	s.SelectorState.Handle(e)
	s.MarkdownState.Handle(e)
	s.PrototypeState.Handle(e)
	s.CanvasState.Handle(e)
	switch e.(type) {
	case EventReloadNotes:
		go s.reloadNotes()
	}
}

func View(s *State) *browser.Node {
	var view *browser.Node
	switch s.SelectorState.SelectedKey {
	case "main", "":
		view = main(s)
	case "markdown":
		view = markdown.View(&s.MarkdownState)
	case "manager":
		view = manager.View(&s.ManagerState)
	case "prototype":
		view = prototype.View(&s.PrototypeState)
	case "canvas":
		view = canvas.View(&s.CanvasState)
	default:
		panic(fmt.Sprintf("unknown key: %s", s.SelectorState.SelectedKey))
	}

	return ui.VStack(
		selector.View(&s.SelectorState, items),
		view,
	)
}

func main(s *State) *browser.Node {
	return ui.VStack(
		s.Theme.Text("Future editions will include a notes app"),
		s.Theme.Button("ReloadNotes").OnClickDispatch(EventReloadNotes{}),
		s.Theme.Text(s.Status),
	)
}

func (s *State) reloadNotes() {
	k := *s.PrivateKey
	s.Status = "loading"
	go browser.Dispatch(nil)
	defer func() { go browser.Dispatch(nil) }()
	c := new(fs.DirServerHTTPClient)

	resp := c.Tree(&fs.DirTreeRequest{
		Public: string(k.Name), Private: k.Private,
		Citizen: ctzn.Name(s.ManagerState.Citizen), Path: fs.Path(s.ManagerState.Path),
		Level: 1, // for now, assume no folders
	})

	if resp.Error != "" {
		s.Status = resp.Error
		return
	}

	s.DirEntries = resp.Entries
	s.Status = ""
}
