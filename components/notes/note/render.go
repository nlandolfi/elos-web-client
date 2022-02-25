package note

import (
	"fmt"
	"io"
	"strings"
)

type renderState struct {
	io.Writer
	indentCount int
	inlining    bool
}

func (s *renderState) indent() {
	for i := 0; i < s.indentCount; i++ {
		fmt.Fprintf(s, "  ")
	}
}

func (s *renderState) printf(format string, a ...interface{}) {
	fmt.Fprintf(s, format, a...)
}

func (s *renderState) openBlock(n int) {
	s.printf(" ")
	for i := 0; i < n; i++ {
		s.printf("{")
	}
	s.nl()
	s.indentCount += 1
}

func (s *renderState) closeBlock(n int) {
	s.indentCount -= 1
	if s.inlining {
		s.nl()
		s.inlining = false
	}
	s.indent()
	for i := 0; i < n; i++ {
		s.printf("}")
	}
	s.nl()
}

func (s *renderState) nl() {
	fmt.Fprint(s, "\n")
}

func (s *renderState) pre(n *Node) {
	switch n.Type {
	case NodeHTML:
		if s.inlining {
			panic("shouldn't be inlining html")
		}
		s.indent()
		s.printf(".html")
		s.openBlock(3)
		for _, l := range n.HTMLInfo.Lines {
			s.printf(l)
			s.nl()
		}
	case NodeDocument:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".doc %s %s %s",
			UnderscoreIfNot(n.DocumentInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.DocumentInfo.Text, " ", "-")),
			UnderscoreIfNot(string(n.DocumentInfo.Type)),
		)
		s.nl()
	case NodeHeader:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.nl()
		s.indent()
		s.printf(".header %d", n.HeaderInfo.Level)
		s.openBlock(1)
	case NodeSection:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.nl()
		s.indent()
		s.printf(".sec %s %s",
			UnderscoreIfNot(n.SectionInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.SectionInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	case NodeParagraph:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".par")
		s.openBlock(1)
	case NodeEquation:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".eq %s", n.EquationInfo.Symbol)
		s.openBlock(1)
	case NodeList:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.nl()
		s.indent()
		switch n.ListInfo.Type {
		case ListOrdered:
			s.printf(".listo")
		default:
			s.printf(".list")
		}
		s.openBlock(1)
	case NodeListItem:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".item")
		s.openBlock(1)
	case NodeExample:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".ex %s %s",
			UnderscoreIfNot(n.ExampleInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.ExampleInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	case NodeDefinition:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".def %s %s",
			UnderscoreIfNot(n.DefinitionInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.DefinitionInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	case NodeCorollary:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".cor %s %s",
			UnderscoreIfNot(n.CorollaryInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.CorollaryInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	case NodeTheorem:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".thm %s %s",
			UnderscoreIfNot(n.TheoremInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.TheoremInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	case NodeImage:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".img %s %s %s",
			UnderscoreIfNot(n.ImageInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.ImageInfo.Text, " ", "-")),
			n.ImageInfo.Path,
		)
		s.openBlock(1)
	case NodeVStack:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".vstack")
		s.openBlock(1)
	case NodeHStack:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".hstack")
		s.openBlock(1)
	case NodeLink:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".link %s %s",
			UnderscoreIfNot(n.LinkInfo.Ref),
			UnderscoreIfNot(strings.ReplaceAll(n.LinkInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	case NodeTex:
		if n.TexInfo.Display {
			if s.inlining {
				s.nl()
			}
			s.inlining = false
			s.indent()
			s.printf(".tex")
			s.openBlock(1)
		} else {
			if !s.inlining {
				s.nl()
				s.indent()
			}
			s.inlining = true
			s.printf("$")
		}
	case NodeComment:
		if s.inlining {
			s.nl()
			s.inlining = false
		}
		s.printf("//")
	case NodeText:
		if !s.inlining {
			s.indent()
		}
		s.inlining = true
		s.printf(n.TextInfo.Text)
	case NodeBold:
		if !s.inlining {
			s.indent()
		}
		s.inlining = true
		s.printf("*")
	case NodeItalics:
		if !s.inlining {
			s.indent()
		}
		s.inlining = true
		s.printf("_")
	case NodeAlgorithm:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".alg %s %s",
			UnderscoreIfNot(n.AlgorithmInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.AlgorithmInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	case NodeProposition:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".prop %s %s",
			UnderscoreIfNot(n.PropositionInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.PropositionInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	case NodeProblem:
		if s.inlining {
			s.nl()
		}
		s.inlining = false
		s.indent()
		s.printf(".prob %s %s",
			UnderscoreIfNot(n.ProblemInfo.Symbol),
			UnderscoreIfNot(strings.ReplaceAll(n.ProblemInfo.Text, " ", "-")),
		)
		s.openBlock(1)
	default:
		panic(fmt.Sprintf("unkown type: %q", n.Type))
	}
}

func UnderscoreIfNot(s string) string {
	if s == "" {
		return "_"
	}
	return s
}

func (s *renderState) post(n *Node) {
	switch n.Type {
	case NodeHTML:
		s.closeBlock(3)
	case NodeDocument:
		s.nl()
	case NodeSection, NodeDefinition, NodeExample, NodeTheorem,
		NodeList, NodeCorollary, NodeParagraph, NodeEquation,
		NodeListItem, NodeImage, NodeVStack, NodeHStack, NodeLink,
		NodeHeader, NodeAlgorithm, NodeProposition, NodeProblem:
		s.closeBlock(1)
	case NodeTex:
		if n.TexInfo.Display {
			s.closeBlock(1)
		} else {
			s.printf("$")
		}
	case NodeText, NodeComment:
	case NodeBold:
		s.printf("*")
	case NodeItalics:
		s.printf("_")
	default:
		panic(fmt.Sprintf("unkown type: %q", n.Type))
	}
}

func Render(w io.Writer, n *Node) {
	var s renderState
	s.Writer = w
	n.Walk(s.pre, s.post)
}
