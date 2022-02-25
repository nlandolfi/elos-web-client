package note

import "fmt"

type NodeType string

const (
	NodeDocument    NodeType = "doc"
	NodeComment     NodeType = "comment"
	NodeHeader      NodeType = "header"
	NodeText        NodeType = "text"
	NodeTex         NodeType = "tex"
	NodeEquation    NodeType = "equation"
	NodeList        NodeType = "list"
	NodeListItem    NodeType = "list-item"
	NodeItalics     NodeType = "italics"
	NodeBold        NodeType = "bold"
	NodeQuote       NodeType = "quote"
	NodeTable       NodeType = "table"
	NodeTableRow    NodeType = "table"
	NodeLink        NodeType = "link"
	NodeRef         NodeType = "ref"
	NodeParagraph   NodeType = "paragraph"
	NodeDefinition  NodeType = "definition"
	NodeTheorem     NodeType = "theorem"
	NodeCorollary   NodeType = "corollary"
	NodeExample     NodeType = "example"
	NodeSection     NodeType = "section"
	NodeImage       NodeType = "image"
	NodeVStack      NodeType = "vstack"
	NodeHStack      NodeType = "hstack"
	NodeTerm        NodeType = "term"
	NodeAlgorithm   NodeType = "algorithm"
	NodeProposition NodeType = "proposition"
	NodeProblem     NodeType = "problem"
	NodeHTML        NodeType = "html"
)

type Node struct {
	Type NodeType

	*DocumentInfo
	*HeaderInfo
	*TextInfo
	*ListInfo
	*EquationInfo
	*TexInfo
	*LinkInfo
	*RefInfo
	*DefinitionInfo
	*CorollaryInfo
	*TheoremInfo
	*PropositionInfo
	*ProblemInfo
	*ExampleInfo
	*SectionInfo
	*ParagraphInfo
	*ImageInfo
	*AlgorithmInfo
	*HTMLInfo

	Children []*Node
}

func (n *Node) Debug() string {
	switch n.Type {
	case NodeDocument:
		return fmt.Sprintf("%s:%s,%s,%s", string(n.Type), n.DocumentInfo.Symbol, n.DocumentInfo.Text, n.DocumentInfo.Type)
	case NodeHeader:
		return fmt.Sprintf("%s:%d", string(n.Type), n.HeaderInfo.Level)
	case NodeItalics:
		return fmt.Sprintf("%s", string(n.Type))
	case NodeList:
		return fmt.Sprintf("%s:%s", string(n.Type), string(n.ListInfo.Type))
	case NodeText:
		return fmt.Sprintf("%s:%s", string(n.Type), string(n.TextInfo.Text))
	case NodeListItem:
		return fmt.Sprintf("%s", string(n.Type))
	default:
		return "no debug info available"
	}
}

type AlgorithmInfo struct {
	Symbol string
	Text   string
}

type PropositionInfo struct {
	Symbol string
	Text   string
}

type ProblemInfo struct {
	Symbol string
	Text   string
}

type ParagraphInfo struct{}

type DocumentInfo struct {
	Symbol string
	Text   string
	Type   DocumentType
}

type DocumentType string

const (
	DocumentArticle = "article"
	DocumentSlides  = "slides"
)

type EquationInfo struct {
	Symbol string
}

type TexInfo struct {
	Display bool
}

type SectionInfo struct {
	Symbol string
	Text   string
}

type HeaderInfo struct {
	Level int
}

type HTMLInfo struct {
	Lines []string
}

type TextInfo struct {
	Text string
}

type ListInfo struct {
	Type ListType
}

type ImageInfo struct {
	Symbol string
	Text   string
	Path   string
}

type ListType string

const (
	ListOrdered   = "ordered"
	ListUnordered = "unordered"
	ListCheck
)

type LinkInfo struct {
	Ref  string
	Text string
}

type RefInfo struct {
	Ref  string
	Text string
}

type DefinitionInfo struct {
	Symbol string
	Text   string
}

type ExampleInfo struct {
	Symbol string
	Text   string
}

type CorollaryInfo struct {
	Symbol string
	Text   string
}

type TheoremInfo struct {
	Symbol string
	Text   string
}

func (root *Node) Walk(pre, post func(n *Node)) {
	walk(root, pre, post)
}

func walk(n *Node, pre, post func(n *Node)) {
	pre(n)
	for _, c := range n.Children {
		walk(c, pre, post)
	}
	post(n)
}

func AssetPaths(n *Node) []string {
	var as []string
	n.Walk(func(n *Node) {
		switch n.Type {
		case NodeImage:
			as = append(as, n.ImageInfo.Path)
		case NodeLink:
			as = append(as, n.LinkInfo.Ref)
		}
	}, func(n *Node) {})
	return as
}

func Items(n *Node) int {
	var i int

	n.Walk(func(c *Node) {
		switch c.Type {
		case NodeListItem:
			i += 1
		}
	}, func(c *Node) {})

	return i
}
