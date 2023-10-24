package webui

import (
	"fmt"
	"go_utils/utils"
)

type Element struct {
	Tag        string
	Attributes map[string]string
	Value      string
	ValueEnd   string
	EndTag     string
	Children   []WebUI
}

func (e *Element) SetID(id string) {
	e.SetAttr("id", id)
}

func (e *Element) SetClass(class string) {
	e.SetAttr("class", class)
}

func (e *Element) SetAttr(k, v string) {
	e.Attributes[k] = v
}

func (e *Element) AddChild(w ...WebUI) {
	for _, v := range w {
		e.Children = append(e.Children, v)
	}
}

func (e *Element) Render() string {
	attrString := ""
	for k, v := range e.Attributes {
		attrString += fmt.Sprintf(" %s=\"%s\"", k, v)
	}
	childString := ""
	for _, v := range e.Children {
		childString += v.Render()
	}
	ret := fmt.Sprintf(
		"<%s%s>%s%s%s%s",
		e.Tag, attrString,
		e.Value, childString, e.ValueEnd,
		e.EndTag,
	)
	return ret
}

func NewElement(tag, value string) *Element {
	return &Element{
		Tag: tag, EndTag: "</" + tag + ">",
		Value:      value,
		Children:   []WebUI{},
		Attributes: map[string]string{},
	}
}

func NewElementWithNoEndTag(tag, value string) *Element {
	return &Element{
		Tag: tag, EndTag: "",
		Value:      value,
		Children:   []WebUI{},
		Attributes: map[string]string{},
	}
}

type GroupElement struct {
	Elements []*Element
}

func NewGroupElement(e ...*Element) *GroupElement {
	return &GroupElement{Elements: e}
}

func (g *GroupElement) Render() string {
	ret := ""
	for _, v := range g.Elements {
		ret += v.Render()
	}
	return ret
}

func NewDiv(w ...WebUI) *Element {
	d := NewElement("div", "")
	d.AddChild(w...)
	return d
}

func NewColumnDiv(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("ccdiv")
	return d
}

func NewDiv2C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-2")
	return d
}

func NewDiv3C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-3")
	return d
}

// level should be h1,h2,h3...
func NewHeader(content, level string) *Element {
	return NewElement(level, content)
}

func NewText(content string) *Element {
	return NewElement("p", content)
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

func NewForm(targetURL, name string, e ...WebUI) *Element {
	f := NewElement("form", "")
	f.Value = fmt.Sprintf("<fieldset><legend>%s</legend>", name)
	f.ValueEnd = "</fieldset>"
	f.SetAttr("action", targetURL)
	f.SetAttr("method", "post")
	f.SetAttr("enctype", "multipart/form-data")
	f.AddChild(e...)
	return f
}

func NewInput(name, inputType, content, id string) *Element {
	i := NewElementWithNoEndTag("input", "")
	i.SetAttr("type", inputType)
	i.SetAttr("name", name)
	i.SetAttr("value", content)
	i.SetAttr("id", id)
	i.SetAttr("autocomplete", "off")
	return i
}

func NewSubmitBtn(name string, id string) *GroupElement {
	s := NewInput(name, "submit", name, id)
	s.SetClass("btn")
	return NewGroupElement(s, NewBR())
}

func NewLabel(content string, target string) *Element {
	l := NewElement("label", content)
	l.SetAttr("for", target)
	return l
}

func NewTextInput(name string) *GroupElement {
	idSuffix := utils.RandomString(5)
	return NewGroupElement(
		NewLabel(name, name+"-"+idSuffix),
		NewInput(name, "text", "", name+"-"+idSuffix), NewBR(),
	)
}

func NewTextInputWithValue(name, value string) *GroupElement {
	idSuffix := utils.RandomString(5)
	return NewGroupElement(
		NewLabel(name, name+"-"+idSuffix),
		NewInput(name, "text", value, name+"-"+idSuffix), NewBR(),
	)
}

func NewCheckBox(name string) *GroupElement {
	idSuffix := utils.RandomString(5)
	return NewGroupElement(
		NewLabel(name, name+"-"+idSuffix),
		NewInput(name, "checkbox", name, name+"-"+idSuffix), NewBR(),
	)
}

func NewRadioInput(name, value string) *GroupElement {
	return NewGroupElement(
		NewInput(name, "radio", value, name+"-"+value),
		NewLabel(value, name+"-"+value), NewBR(),
	)
}

func NewTable(w ...WebUI) *Element {
	t := NewElement("table", "")
	t.Value = "<tbody>"
	t.ValueEnd = "</tbody>"
	t.AddChild(w...)
	return t
}

func NewTableRow(header bool, data ...string) *Element {
	rowType := "td"
	if header {
		rowType = "th"
	}
	tr := NewElement("tr", "")
	for _, v := range data {
		th := NewElement(rowType, v)
		tr.AddChild(th)
	}
	return tr
}
