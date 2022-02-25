package note

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func MustParseFromFile(s string) *Node {
	f, err := os.Open(s)
	if err != nil {
		log.Fatalf("md.MustParseFromFile os.Open %v", err)
	}
	defer f.Close()
	n, err := Parse(f)
	if err != nil {
		log.Fatalf("md.MustParseFromFile  Parse %v", err)
	}
	return n
}

func MustParse(r io.Reader) *Node {
	n, err := Parse(r)
	if err != nil {
		log.Fatalf("md.MustParse: %v", err)
	}
	return n
}

func Parse(r io.Reader) (*Node, error) {
	var p parseState
	p.Scanner = bufio.NewScanner(r)
	p.stack = []*Node{{}}
	if err := p.run(); err != nil {
		return nil, fmt.Errorf("md.Parse: %v", err)
	}
	if len(p.stack[0].Children) > 0 {
		return p.stack[0].Children[0], nil
	} else {
		return Document("", "", DocumentArticle), nil
	}
}

type parseState struct {
	*bufio.Scanner
	stack      []*Node
	lineNumber int
	inTex      bool
	inHTML     bool
}

func (s *parseState) addchild(n *Node) {
	last := s.stack[len(s.stack)-1]
	last.Children = append(last.Children, n)
}

func (s *parseState) push(n *Node) {
	if len(s.stack) > 0 {
		s.addchild(n)
	}

	s.stack = append(s.stack, n)
}

func (s *parseState) pop() {
	s.stack = s.stack[:len(s.stack)-1]
}

func (s *parseState) current() *NodeType {
	if len(s.stack) == 0 {
		return nil
	}
	return &s.stack[len(s.stack)-1].Type
}

type lineMode string

const (
	lineTex     lineMode = "tex"
	lineBold             = "bold"
	lineItalics          = "italics"
)

type lineState struct {
	children []*Node
	mode     lineMode
	b        strings.Builder
}

func modeFor(r rune) lineMode {
	switch r {
	case '$':
		return lineTex
	case '*':
		return lineBold
	case '_':
		return lineItalics
	default:
		return ""
	}
}

func (s *lineState) node() *Node {
	switch s.mode {
	case lineTex:
		return &Node{
			Type: NodeTex,
			TexInfo: &TexInfo{
				Display: false,
			},
			Children: []*Node{Text(s.b.String())},
		}
	case lineBold:
		return &Node{
			Type:     NodeBold,
			Children: []*Node{Text(s.b.String())},
		}
	case lineItalics:
		return &Node{
			Type:     NodeItalics,
			Children: []*Node{Text(s.b.String())},
		}
	case "":
		return Text(s.b.String())
	default:
		panic(fmt.Sprintf("lineState.mode unknown mode: %s", s.mode))
	}
}

func (s *lineState) reset() {
	s.mode = ""
	s.b.Reset()
}

func (s *parseState) consumeText(text string) error {
	var ls lineState

	for i, r := range text {
		switch m := modeFor(r); m {
		case lineTex, lineBold, lineItalics:
			if m == lineTex && i+1 < len(text) && modeFor(rune(text[i+1])) == lineTex {
				// skip this guy as it is an opening or closing $ for display mode.
				continue
			}
			// if we are alread in this mode, end it.
			if ls.mode == m {
				ls.children = append(ls.children, ls.node())
				ls.reset()
				continue
			}

			// otherwise start it, unless we are in tex mode
			if ls.mode == lineTex {
				ls.b.WriteRune(r)
				continue
			}
			ls.children = append(ls.children, ls.node())
			ls.reset()
			ls.mode = m
		default:
			ls.b.WriteRune(r)
		}
	}

	if ls.mode != "" {
		return fmt.Errorf("line %q \n unfinished mode %q", text, ls.mode)
	}

	ls.children = append(ls.children, ls.node())

	for _, c := range ls.children {
		s.addchild(c)
	}

	return nil
}

func (s *parseState) consumeLine(line string) error {
	s.lineNumber += 1
	if line == "" {
		return nil
	}
	if s.inHTML {
		if strings.TrimSpace(line) == "}}}" {
			s.inHTML = false
			s.pop()
			return nil
		}
		n := s.stack[len(s.stack)-1].HTMLInfo
		n.Lines = append(n.Lines, line)
		return nil
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil
	}
	if line == "}" {
		s.pop()
		return nil
	}
	if line == "$$" {
		if s.inTex {
			s.pop()
			return nil
		}
		s.push(&Node{
			Type: NodeTex,
			TexInfo: &TexInfo{
				Display: true,
			},
		})
		return nil
	}

	if len(line) >= 2 && line[:2] == "//" {
		s.addchild(&Node{
			Type:     NodeComment,
			Children: []*Node{Text(line[2:])},
		})
		return nil
	}

	if line[0] == '.' {
		args := strings.Fields(line[1:])
		if len(args) == 0 {
			// return fmt.Errorf("line %d, just a line with %q?", s.lineNumber, line)
			s.addchild(Text(line))
			return nil
		}
		switch cmd := args[0]; cmd {
		case "html":
			s.inHTML = true
			s.push(&Node{
				Type:     NodeHTML,
				HTMLInfo: &HTMLInfo{},
			})
		case "doc":
			if len(args) == 3 {
				args = append(args, "article")
			}
			s.push(&Node{
				Type: NodeDocument,
				DocumentInfo: &DocumentInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
					Type:   DocumentType(args[3]),
				},
			})
		case "header":
			i, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			s.push(&Node{
				Type: NodeHeader,
				HeaderInfo: &HeaderInfo{
					Level: int(i),
				},
			})
		case "sec":
			s.push(&Node{
				Type: NodeSection,
				SectionInfo: &SectionInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
				},
			})
		case "p", "par":
			s.push(&Node{
				Type: NodeParagraph,
			})
		case "eq":
			s.push(&Node{
				Type: NodeEquation,
				EquationInfo: &EquationInfo{
					Symbol: args[1],
				},
			})
		case "tex":
			s.push(&Node{
				Type: NodeTex,
				TexInfo: &TexInfo{
					Display: true,
				},
			})
		case "list":
			s.push(&Node{
				Type: NodeList,
				ListInfo: &ListInfo{
					Type: ListUnordered,
				},
			})
		case "listo":
			s.push(&Node{
				Type: NodeList,
				ListInfo: &ListInfo{
					Type: ListOrdered,
				},
			})
		case "item":
			s.push(&Node{
				Type: NodeListItem,
			})
		case "ex":
			s.push(&Node{
				Type: NodeExample,
				ExampleInfo: &ExampleInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
				},
			})
		case "def":
			s.push(&Node{
				Type: NodeDefinition,
				DefinitionInfo: &DefinitionInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
				},
			})
		case "cor":
			s.push(&Node{
				Type: NodeCorollary,
				CorollaryInfo: &CorollaryInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
				},
			})
		case "thm":
			s.push(&Node{
				Type: NodeTheorem,
				TheoremInfo: &TheoremInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
				},
			})
		case "img", "image":
			if len(args) != 4 {
				return fmt.Errorf("error parsing img")
			}
			s.push(&Node{
				Type: NodeImage,
				ImageInfo: &ImageInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
					Path:   args[3],
				},
			})
		case "vstack":
			s.push(&Node{
				Type: NodeVStack,
			})
		case "hstack":
			s.push(&Node{
				Type: NodeHStack,
			})

		case "link":
			if len(args) != 3 {
				return fmt.Errorf("error parsing link")
			}
			s.push(&Node{
				Type: NodeLink,
				LinkInfo: &LinkInfo{
					Ref: args[1],
					// this text appears ot be a label?
					Text: strings.ReplaceAll(args[2], "-", " "),
				},
			})
		case "alg":
			s.push(&Node{
				Type: NodeAlgorithm,
				AlgorithmInfo: &AlgorithmInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
				},
			})
		case "prop":
			s.push(&Node{
				Type: NodeProposition,
				PropositionInfo: &PropositionInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
				},
			})
		case "prob":
			s.push(&Node{
				Type: NodeProblem,
				ProblemInfo: &ProblemInfo{
					Symbol: args[1],
					Text:   strings.ReplaceAll(args[2], "-", " "),
				},
			})
		default:
			return fmt.Errorf("line %d: unknown command: %q", s.lineNumber, cmd)
		}

		return nil
	}

	if s.current() != nil && (*s.current() == NodeTex || *s.current() == NodeEquation) {
		s.addchild(Text(line))
		return nil
	}

	if err := s.consumeText(line); err != nil {
		return fmt.Errorf("line %d: %s", s.lineNumber, err)
	}

	return nil
}

func (s *parseState) run() error {
	for s.Scan() {
		line := s.Text()
		if err := s.consumeLine(line); err != nil {
			return err
		}
	}
	if err := s.Err(); err != nil {
		return err
	}
	return nil
}
