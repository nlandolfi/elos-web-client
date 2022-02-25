package note

func Document(symbol, text string, t DocumentType, children ...*Node) *Node {
	return &Node{
		Type: NodeDocument,
		DocumentInfo: &DocumentInfo{
			Symbol: symbol,
			Text:   text,
			Type:   t,
		},
		Children: children,
	}
}

func Header(level int, children ...*Node) *Node {
	return &Node{
		Type:       NodeHeader,
		HeaderInfo: &HeaderInfo{Level: level},
		Children:   children,
	}
}

func Text(text string) *Node {
	return &Node{
		Type: NodeText,
		TextInfo: &TextInfo{
			Text: text,
		},
	}
}

func TexInline(tex string) *Node {
	return &Node{
		Type: NodeTex,
		TexInfo: &TexInfo{
			Display: false,
		},
		Children: []*Node{Text(tex)},
	}
}

func TexDisplay(tex string) *Node {
	return &Node{
		Type: NodeTex,
		TexInfo: &TexInfo{
			Display: true,
		},
		Children: []*Node{Text(tex)},
	}
}

func OrderedList(items []string) *Node {
	children := make([]*Node, len(items))
	for i, s := range items {
		children[i] = &Node{
			Type: NodeListItem,
			Children: []*Node{
				Text(s),
			},
		}
	}
	return &Node{
		Type: NodeList,
		ListInfo: &ListInfo{
			Type: ListOrdered,
		},
		Children: children,
	}
}

func ListItem(children ...*Node) *Node {
	return &Node{
		Type:     NodeListItem,
		Children: children,
	}
}

func UnorderedList(items []string) *Node {
	children := make([]*Node, len(items))
	for i, s := range items {
		children[i] = &Node{
			Type: NodeListItem,
			Children: []*Node{
				Text(s),
			},
		}
	}
	return &Node{
		Type: NodeList,
		ListInfo: &ListInfo{
			Type: ListUnordered,
		},
		Children: children,
	}
}

func Italics(text string) *Node {
	return &Node{
		Type: NodeItalics,
		Children: []*Node{
			Text(text),
		},
	}
}

func Bold(text string) *Node {
	return &Node{
		Type: NodeBold,
		Children: []*Node{
			Text(text),
		},
	}
}

func Quote(text string) *Node {
	return &Node{
		Type: NodeBold,
		Children: []*Node{
			Text(text),
		},
	}
}

func Section(symbol, text string, children ...*Node) *Node {
	return &Node{
		Type: NodeSection,
		SectionInfo: &SectionInfo{
			Symbol: symbol,
			Text:   text,
		},
		Children: children,
	}
}

func Image(symbol, text, path string) *Node {
	return &Node{
		Type: NodeImage,
		ImageInfo: &ImageInfo{
			Symbol: symbol,
			Text:   text,
			Path:   path,
		},
	}
}

func VStack(children ...*Node) *Node {
	return &Node{
		Type:     NodeVStack,
		Children: children,
	}
}

func HStack(children ...*Node) *Node {
	return &Node{
		Type:     NodeHStack,
		Children: children,
	}
}

func Paragraph(text string) *Node {
	return &Node{
		Type: NodeParagraph,
		Children: []*Node{
			Text(text),
		},
	}
}
