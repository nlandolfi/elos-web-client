package editor

import (
	"log"

	"github.com/nlandolfi/spin/infra/ctzn"
	"github.com/nlandolfi/spin/infra/fs"
	"github.com/nlandolfi/spin/infra/key"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	Theme      *ui.Theme        `json:"-"`
	PrivateKey **key.PrivateKey `json:"-"`

	File File
}

func (s *State) Handle(e browser.Event) {
	switch e.(type) {
	case EventReload:
		log.Print("reload")
		go s.File.Reload()
	case EventSave:
		log.Print("save")
		go s.File.Save()
	}
}

type EventReload struct{}
type EventSave struct{}

func (s *State) Rewire(t *ui.Theme, k **key.PrivateKey) {
	s.Theme = t
	s.PrivateKey = k
	s.File.PrivateKey = k
}

func View(s *State) *browser.Node {
	return ui.VStack(
		ui.HStack(
			s.Theme.TextInput(&s.File.Citizen).Placeholder("citizen").FlexGrow("0.01"),
			s.Theme.TextInput(&s.File.Path).Placeholder("path to spin file...").FlexGrow("1").
				OnKeyDown(func(e dom.Event) {
					if e.KeyCode() == 13 { // enter
						go browser.Dispatch(EventReload{}) // todo change name of action
					}
				}),
			s.Theme.Text(s.File.Status),
			s.Theme.Button("Reload").OnClick(browser.Dispatcher(EventReload{})),
			s.Theme.Button("Save").OnClick(browser.Dispatcher(EventSave{})),
		),
		s.Theme.TextArea(&s.File.Text).MinHeight(browser.Size{Value: 500, Unit: browser.UnitPX}),
	).PaddingPX(10)
}

func (s *State) Open(c ctzn.Name, p fs.Path) {
	s.File.Citizen = string(c)
	s.File.Path = string(p)
	s.File.Reload()
}

/*
func loadText(s *State) {
	k := *s.PrivateKey

	s.Status = "loading..."
	c := new(txt.TxtServerHTTPClient)

	resp := c.Reload(&txt.TxtReloadRequest{
		Public: string(k.Name), Private: k.Private,
		Citizen: k.Citizen,
		Path:    fs.Path(s.Path),
	})

	if resp.Error != "" {
		s.Error = resp.Error
		return
	}

	s.Text = resp.Snapshot
	s.Shadow = resp.Snapshot
	s.ShadowSequence = resp.Sequence
	s.Error = ""
}

func save(s *State) {
	k := *s.PrivateKey
	s.Error = "saving..."

	// these are the ops we've seen so far
	sentSnapshot := s.Text
	ops := txtops.DiffOps(txtops.Diffs(s.Shadow, s.Text))

	c := new(txt.TxtServerHTTPClient)

	resp := c.Commit(&txt.TxtCommitRequest{
		Public: string(k.Name), Private: k.Private,
		Citizen:  k.Citizen,
		Path:     fs.Path(s.Path),
		Sequence: s.ShadowSequence,
		Ops:      ops,
	})

	if resp.Error != "" {
		s.Error = resp.Error
		return
	}

	// these are the ops we've seen since the commit request
	sinceSaveOps := txtops.DiffOps(txtops.Diffs(sentSnapshot, s.Text))

	// now we update the shadow
	s.Shadow = resp.Snapshot
	s.ShadowSequence = resp.Sequence

	for _, serverOp := range resp.Ops {
		for _, localOp := range sinceSaveOps {
			txtops.DiffOpAdjust(localOp, &serverOp)
		}
	}

	s.Text = s.Shadow
	for _, op := range sinceSaveOps {
		s.Text = txtops.DiffOpApply(s.Text, op)
	}

	s.Error = ""

}
*/
