package webui

// level should be h1,h2,h3...
func NewHeader(content, level string) *Element {
	return NewElement(level, content)
}

func NewText(content string) *Element {
	return NewElement("p", content)
}

func NewPreText(content string) *Element {
	return NewElement("pre", content)
}

func NewLink(content string, link string) *Element {
	e := NewElement("a", content)
	e.SetAttr("href", link)
	return e
}

func NewLinkBtn(content string, link string) *Element {
	e := NewElement("a", content)
	e.SetAttr("href", link)
	e.SetClass("btn")
	return e
}

func NewButton(content, link string) *Element {
	e := NewElement("a", content)
	e.SetAttr("href", link)
	b := NewElement("button", "")
	b.AddChild(e)
	return b
}

func NewBR() *Element {
	return NewElementWithNoEndTag("br", "")
}
