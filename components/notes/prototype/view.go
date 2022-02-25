package prototype

import (
	"bytes"
	"fmt"
	"log"

	"github.com/nlandolfi/elos/web-client/components/notes/note"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"github.com/spinsrv/browser/ui"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type State struct {
	Theme *ui.Theme

	Model     string
	Debugging bool
	Raw       string
	Root      *note.Node `json:"-"`
	Selected  *note.Node `json:"-"`

	Status string

	Selection dom.Selection `json:"-"`
}

func (s *State) Handle(e browser.Event) {
	switch e := e.(type) {
	case EventKey:
		b := bytes.NewBufferString(s.Raw)
		n, err := note.Parse(b)
		if err != nil {
			s.Status = err.Error()
			return // don't worry, may be mid typing
			//			log.Fatal(err)
		} else {
			s.Status = ""
			s.Root = n
		}
		/*
			var bb bytes.Buffer
			note.Render(&bb, s.Root)
			s.Raw = bb.String()
		*/
		// overwrites the cursor position :/

	case EventFormat:
		b := bytes.NewBufferString(s.Raw)
		n, err := note.Parse(b)
		if err != nil {
			s.Status = err.Error()
			return // don't worry, may be mid typing
			//			log.Fatal(err)
		} else {
			s.Status = ""
			s.Root = n
		}
		var bb bytes.Buffer
		note.Render(&bb, s.Root)
		s.Raw = bb.String()
		// overwrites the cursor position :/
	case EventToggleDebugging:
		s.Debugging = !s.Debugging
	case EventAppendNode:
		if s.Selected != nil {
			s.Selected.Children = append(s.Root.Children, &e.Node)
		} else {
			s.Root.Children = append(s.Root.Children, &e.Node)
		}
		var b bytes.Buffer
		note.Render(&b, s.Root)
		s.Raw = b.String()
	case EventReset:
		s.Root = note.Document("", "", note.DocumentArticle)
	case EventSelectNode:
		if s.Selected == e.Node {
			s.Selected = nil
		} else {
			s.Selected = e.Node
		}
	}
}

type EventToggleRaw struct{}
type EventAppendNode struct{ Node note.Node }
type EventReset struct{}
type EventSelectNode struct{ Node *note.Node }
type EventToggleDebugging struct{ Node *note.Node }
type EventKey struct{}
type EventFormat struct{}

func View(s *State) *browser.Node {
	return ui.VStack(
		ui.HStack(
			s.Theme.Button("toggle").OnClickDispatch(EventToggleRaw{}),
			s.Theme.Button("Format").OnClickDispatch(EventFormat{}),
			s.Theme.Button("debug").OnClickDispatch(EventToggleDebugging{}),
			s.Theme.Button("reset").OnClickDispatch(EventReset{}),
			s.Theme.Textf("AnchorOffset: %d; FocusOffset: %d; IsCollapsed: %t; RangeCount: %d; Type: %s", s.Selection.AnchorOffset(), s.Selection.FocusOffset(), s.Selection.IsCollapsed(), s.Selection.RangeCount(), s.Selection.Type()),
		),
		s.Theme.Textf("AchorNode: %+v", s.Selection.AnchorNode()),
		s.Theme.Text(s.Status),
		ui.HStack(
			ui.VStack(
				ui.HStack(
					s.Theme.Button("H1").
						OnClickDispatch(EventAppendNode{Node: *note.Header(1)}),
					s.Theme.Button("H2").
						OnClickDispatch(EventAppendNode{Node: *note.Header(2)}),
					s.Theme.Button("H3").
						OnClickDispatch(EventAppendNode{Node: *note.Header(3)}),
					s.Theme.Button("OL").
						OnClickDispatch(EventAppendNode{Node: *note.OrderedList(nil)}),
					s.Theme.Button("UL").
						OnClickDispatch(EventAppendNode{Node: *note.UnorderedList(nil)}),
					s.Theme.Button("LI").
						OnClickDispatch(EventAppendNode{Node: *note.ListItem()}),
				),
				s.Theme.Card(s.render(s.Root)).WidthPX(500).
					OnClickDispatch(nil),
			).WidthPG(50),
			s.Theme.TextArea(&s.Raw).HeightPX(400).FontFamily("monospace").WidthPG(50),
		),
		ui.OnlyIf(s.Debugging,
			func() *browser.Node {
				var b bytes.Buffer
				var level int

				s.Root.Walk(func(n *note.Node) {
					for i := 0; i < level; i++ {
						fmt.Fprintf(&b, "  ")
					}

					fmt.Fprintf(&b, "%s\n", n.Debug())
					level++
				}, func(n *note.Node) {
					level--
				})

				out := b.String()
				return s.Theme.TextArea(&out).FontFamily("monospace").HeightPX(800)
			},
		),
	).OnKeyUp(func(e dom.Event) {
		go browser.Dispatch(EventKey{})
	})
}

func (s *State) block(n *browser.Node, nn *note.Node, selected bool) *browser.Node {
	var color = "gray"
	if selected {
		color = "red"
	}
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Span,
		Style: browser.Style{
			Border: browser.Border{
				Type: browser.BorderSolid,
				Width: browser.Size{
					Value: 1,
					Unit:  browser.UnitPX,
				},
				Color: color,
			},
		},
		Children: []*browser.Node{
			ui.VStack(
				ui.HStack(
					ui.VStack(
						ui.If(selected,
							func() *browser.Node {
								return s.Theme.Button("S")
							},
							func() *browser.Node {
								return s.Theme.Button("U")
							},
						).OnClickDispatch(EventSelectNode{nn}),
						s.Theme.Button("D"),
					),
					n,
				),
				s.Theme.Text(nn.Debug()).FontFamily("monospace").FontSizePX(12),
			),
		},
	}
}

func (s *State) render(n *note.Node) *browser.Node {
	root := &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Span,
	}

	var stack []*browser.Node = []*browser.Node{root}

	n.Walk(func(n *note.Node) {

		//		log.Printf("Walk n.Type = %s; entering=%t: stack: %v", n.Type, true, stack)

		parent := stack[len(stack)-1]
		var nextNode *browser.Node

		switch n.Type {
		case note.NodeBold:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.B,
			}
		case note.NodeDocument:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Span,
				Attr: []*html.Attribute{
					&html.Attribute{
						Key: "contenteditable",
						Val: "true",
					},
				},
			}
		case note.NodeHeader:
			switch n.HeaderInfo.Level {
			case 1:
				nextNode = s.Theme.H1()
			case 2:
				nextNode = s.Theme.H2()
			case 3:
				nextNode = s.Theme.H3()
			default:
				log.Fatalf("unhandled header level: %d", n.HeaderInfo.Level)
			}
		case note.NodeImage:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Img,
				Attr: []*html.Attribute{
					&html.Attribute{
						Key: "src",
						Val: n.ImageInfo.Path,
					},
					&html.Attribute{
						Key: "alt",
						Val: n.ImageInfo.Text,
					},
				},
			}
		case note.NodeItalics:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.I,
			}
		case note.NodeLink:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Img,
				Attr: []*html.Attribute{
					&html.Attribute{
						Key: "href",
						Val: n.LinkInfo.Ref,
					},
				},
			}
		case note.NodeList:
			switch n.ListInfo.Type {
			case note.ListOrdered:
				nextNode = &browser.Node{
					Type:     html.ElementNode,
					DataAtom: atom.Ol,
				}
			case note.ListUnordered:
				nextNode = &browser.Node{
					Type:     html.ElementNode,
					DataAtom: atom.Ul,
				}
			default:
				panic(fmt.Sprintf("unhandled list type: %s", n.ListInfo.Type))
			}
		case note.NodeListItem:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Li,
			}
		case note.NodeText:
			nextNode = s.Theme.Text(n.TextInfo.Text)
		default:
			log.Fatalf("unhandled node type: %s", n.Type)
		}

		parent.Children = append(parent.Children, nextNode)
		if isContainer(n) {
			stack = append(stack, nextNode)
		}

		return
	}, func(n *note.Node) {
		//		log.Printf("Walk n.Type = %s; entering=%t: stack: %v", n.Type, false, stack)
		if isContainer(n) {
			stack = stack[:len(stack)-1]
			return
		}
	})

	if len(stack) != 1 {
		for i := len(stack) - 1; i >= 0; i-- {
			log.Printf("stack[%d].Type = %s", i, stack[i].DataAtom)
		}
		log.Fatalf("len(stack) got %d, want %d", len(stack), 1)
	}

	return stack[0]
}

func isContainer(n *note.Node) bool {
	switch n.Type {
	case note.NodeDocument,
		note.NodeComment,
		note.NodeHeader,
		note.NodeList,
		note.NodeListItem,
		note.NodeItalics,
		note.NodeBold, note.NodeLink: // incomplete
		return true
	default:
		return false
	}
}
