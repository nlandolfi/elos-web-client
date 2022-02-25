package text

import (
	"fmt"
	"log"

	"github.com/nlandolfi/spin/infra/ctzn"
	sfs "github.com/nlandolfi/spin/infra/fs"
	"github.com/nlandolfi/spin/infra/key"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
)

type State struct {
	Path string
	Text string

	Error string
}

func View(k *key.PrivateKey, s *State) *browser.Node {
	return ui.VStack(
		ui.HStack(
			ui.TextInput(&s.Path).Placeholder("path to spin file...").Styled(&browser.Style{FlexGrow: "1"}),
			ui.Button("Reload").OnClick(func(e dom.Event) {
				go loadText(k, s)
			}),
			ui.Button("Save").OnClick(func(e dom.Event) {
				go save(k, s)
			}),
		),
		ui.Text(s.Error),
		ui.TextArea(&s.Text).Styled(&browser.Style{
			MinHeight: &browser.Size{
				Value: 20,
				Unit:  browser.UnitEM,
			},
		}),
	).Padding(10)
}

func fs(pu, pr string, c ctzn.Name) sfs.System {
	store := &sfs.StoreServerStore{
		Public:      pu,
		Private:     pr,
		StoreServer: new(sfs.StoreServerHTTPClient),
	}
	dir := &sfs.DirServerDir{
		Public:    pu,
		Private:   pr,
		DirServer: new(sfs.DirServerHTTPClient),
	}
	return sfs.NewSystem(c, dir, store)
}

func loadText(k *key.PrivateKey, s *State) {
	s.Error = "loading..."

	files := fs(string(k.Name), k.Private, k.Citizen)

	bs, err := files.ReadFile(sfs.Path(s.Path))
	if err != nil {
		log.Print("couldn't open")
		s.Error = err.Error()
	}

	s.Text = string(bs)
	s.Error = ""
}

func save(k *key.PrivateKey, s *State) {
	s.Error = "saving..."

	files := fs(string(k.Name), k.Private, k.Citizen)

	f, err := files.Open(sfs.Path(s.Path))
	if err != nil {
		s.Error = err.Error()
	}

	f.Truncate()

	fmt.Fprint(f, s.Text)

	if err := f.Close(); err != nil {
		s.Error = err.Error()
	}

	s.Error = ""
}
