package markdown

import (
	"log"
	"strings"

	"github.com/russross/blackfriday/v2"
	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type State struct {
	Theme *ui.Theme

	Markdown string
}

func (s *State) Handle(e browser.Event) {
	switch e.(type) {
	case EventFormatMarkdown:
		log.Print("formatting!")
		formatted, err := Format([]byte(s.Markdown), nil)
		if err != nil {
			log.Fatal(err)
		}
		s.Markdown = string(formatted)
		s.Markdown = strings.Replace(s.Markdown, "\t", " ", -1)
	}
}

type EventFormatMarkdown struct{}

func View(s *State) *browser.Node {
	bf := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	n := bf.Parse([]byte(s.Markdown))

	return ui.VStack(
		ui.Button("Format").OnClickDispatch(EventFormatMarkdown{}),
		ui.HStack(
			ui.Div(
				ui.VStack(
					s.Theme.Text("Editor..."),
					s.Theme.TextArea(&s.Markdown).HeightPX(800).FontFamily("monospace"),
				),
			).FlexGrow("1").Width(browser.Size{Value: 50, Unit: browser.UnitPG}),
			ui.VStack(
				s.Theme.Text("Rendered..."),
				s.Theme.Card(render(s.Theme, n)).
					MinHeight(browser.Size{Value: 3, Unit: browser.UnitEM}).
					//		MaxWidth(browser.Size{Value: 50, Unit: browser.UnitPG}).
					PaddingPX(10),
			).FlexGrow("1").Width(browser.Size{Value: 50, Unit: browser.UnitPG}).FontFamily("monospace"),
		).PaddingPX(20),
	)
}

func render(th *ui.Theme, n *blackfriday.Node) *browser.Node {
	root := &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Span,
	}

	var stack []*browser.Node = []*browser.Node{root}

	n.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		// always just exit
		if !entering {
			log.Printf("Walk n.Type = %s; entering=%t: stack: %v", n.Type, entering, stack)
			if isContainer(n) {
				stack = stack[:len(stack)-1]
				return blackfriday.GoToNext
			}
		}

		log.Printf("Walk n.Type = %s; entering=%t: stack: %v", n.Type, entering, stack)

		parent := stack[len(stack)-1]
		var nextNode *browser.Node

		switch n.Type {
		case blackfriday.BlockQuote:
			// TODO: is this right? - NCL
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Blockquote,
			}
		case blackfriday.Code, blackfriday.CodeBlock:
			// TODO: is this right? - NCL
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Code,
				Children: []*browser.Node{
					&browser.Node{
						Type: html.TextNode,
						Data: string(n.Literal),
					},
				},
			}
		case blackfriday.Del:
			// TODO: is this right? - NCL
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Code,
			}
		case blackfriday.Document:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Span,
			}
		case blackfriday.Emph:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.I,
			}
		case blackfriday.HTMLBlock:
			log.Fatal("html block block not implemented")
		case blackfriday.HTMLSpan:
			log.Fatal("html span block not implemented")
		case blackfriday.Hardbreak:
			// TODO: is this right? - NCL 2/15/22
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Br,
			}
		case blackfriday.Heading:
			switch n.HeadingData.Level {
			case 1:
				nextNode = th.H1()
			case 2:
				nextNode = th.H2()
			case 3:
				nextNode = th.H3()
			default:
				log.Fatalf("heading level %d too high", n.HeadingData.Level)
			}
		case blackfriday.HorizontalRule:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Hr,
			}
		case blackfriday.Image:
			// TODO: complete - NCL 2/15/22
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Img,
			}
		case blackfriday.Item:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Li,
			}
		case blackfriday.Link:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.A,
				Attr: []*html.Attribute{
					&html.Attribute{
						Key: atom.Href.String(),
						Val: string(n.LinkData.Destination),
					},
					&html.Attribute{
						Key: atom.Title.String(),
						Val: string(n.LinkData.Title),
					},
				},
				//				Children: []*browser.Node{th.Text(string(n.LinkData.Title))},
			}
		case blackfriday.List:
			if n.ListData.ListFlags&blackfriday.ListTypeOrdered != 0 {
				nextNode = &browser.Node{
					Type:     html.ElementNode,
					DataAtom: atom.Ol,
				}
			} else {
				nextNode = &browser.Node{
					Type:     html.ElementNode,
					DataAtom: atom.Ul,
				}
			}
		case blackfriday.Paragraph:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.P,
			}
		case blackfriday.Softbreak:
			log.Fatal("softbreak not implemented")
		case blackfriday.Strong:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.B,
			}
		case blackfriday.Table:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Table,
			}
		case blackfriday.TableBody:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Tbody,
			}
		case blackfriday.TableCell:
			nextNode = (&browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Td,
			}).Border(browser.Border{
				Width: browser.Size{Value: 1, Unit: browser.UnitPX},
				Type:  browser.BorderSolid,
				Color: th.TextColor,
			})
		case blackfriday.TableHead:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Thead,
			}
		case blackfriday.TableRow:
			nextNode = &browser.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Tr,
			}
		case blackfriday.Text:
			nextNode = &browser.Node{
				Type: html.TextNode,
				Data: string(n.Literal),
			}
		default:
			log.Fatalf("unhandled node type: %s", n.Type)
		}

		parent.Children = append(parent.Children, nextNode)
		if isContainer(n) {
			stack = append(stack, nextNode)
		}

		return blackfriday.GoToNext
	})

	if len(stack) != 1 {
		for i := len(stack) - 1; i >= 0; i-- {
			log.Printf("stack[%d].Type = %s", i, stack[i].DataAtom)
		}
		log.Fatalf("len(stack) got %d, want %d", len(stack), 1)
	}

	return stack[0]
}

func isContainer(n *blackfriday.Node) bool {
	switch n.Type {
	case blackfriday.Document,
		blackfriday.BlockQuote,
		blackfriday.List,
		blackfriday.Item,
		blackfriday.Paragraph,
		blackfriday.Heading,
		blackfriday.Emph,
		blackfriday.Strong,
		blackfriday.Del,
		blackfriday.Link,
		blackfriday.Image,
		blackfriday.Table,
		blackfriday.TableHead,
		blackfriday.TableBody,
		blackfriday.TableRow,
		blackfriday.TableCell:
		return true
	default:
		return false
	}
}
