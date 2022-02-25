package editor

import (
	"time"

	"github.com/nlandolfi/spin/apps/txt"
	"github.com/nlandolfi/spin/infra/ctzn"
	"github.com/nlandolfi/spin/infra/fs"
	"github.com/nlandolfi/spin/infra/key"
	"github.com/nlandolfi/spin/infra/txtops"
	uuid "github.com/satori/go.uuid"
	"github.com/spinsrv/browser"
)

type File struct {
	PrivateKey     **key.PrivateKey
	Citizen        string
	Path           string
	ShadowSequence int
	Shadow         string
	Text           string
	Status         string
}

func New(k **key.PrivateKey, p string) *File {
	return &File{
		PrivateKey: k,
		Path:       p,
	}
}

func (f *File) Reload() {
	k := *f.PrivateKey

	f.Status = "loading..."
	go browser.Dispatch(nil)
	defer func() { go browser.Dispatch(nil) }()
	c := new(txt.TxtServerHTTPClient)

	resp := c.Reload(&txt.TxtReloadRequest{
		Public: string(k.Name), Private: k.Private,
		Citizen: ctzn.Name(f.Citizen),
		Path:    fs.Path(f.Path),
	})

	if resp.Error != "" {
		f.Status = resp.Error
		return
	}

	f.Text = resp.Snapshot
	f.Shadow = resp.Snapshot
	f.ShadowSequence = resp.Sequence
	f.Status = ""
}

func (f *File) Save() {
	k := *f.PrivateKey
	defer func() { go browser.Dispatch(nil) }()
	f.Status = "saving..."
	go browser.Dispatch(nil) // TODO

	// these are the ops we've seen so far
	sentSnapshot := f.Text
	ops := txtops.DiffOps(txtops.Diffs(f.Shadow, f.Text))
	for _, op := range ops {
		op.Citizen = k.Citizen
		op.Time = time.Now()
		op.ID = uuid.NewV4().String()
	}

	c := new(txt.TxtServerHTTPClient)

	resp := c.Commit(&txt.TxtCommitRequest{
		Public: string(k.Name), Private: k.Private,
		Citizen:  ctzn.Name(f.Citizen),
		Path:     fs.Path(f.Path),
		Sequence: f.ShadowSequence,
		Ops:      ops,
	})

	if resp.Error != "" {
		f.Status = resp.Error
		return
	}

	// these are the ops we've seen since the commit request
	sinceSaveOps := txtops.DiffOps(txtops.Diffs(sentSnapshot, f.Text))

	// now we update the shadow
	f.Shadow = resp.Snapshot
	f.ShadowSequence = resp.Sequence

	for _, serverOp := range resp.Ops {
		for _, localOp := range sinceSaveOps {
			txtops.DiffOpAdjust(localOp, &serverOp)
		}
	}

	f.Text = f.Shadow
	for _, op := range sinceSaveOps {
		f.Text = txtops.DiffOpApply(f.Text, op)
	}

	f.Status = ""
}
