package webui

import (
	"encoding/base64"
	"fmt"
	"go_utils/utils"
	"os"
)

type Element struct {
	Tag        string
	Attributes map[string]string
	Value      string
	ValueEnd   string
	EndTag     string
	Style      map[string]string
	Children   []WebUI
}

type GroupElement struct {
	Elements []*Element
}

func (e *Element) SetAttr(k, v string) {
	e.Attributes[k] = v
}

func (e *Element) SetID(id string) {
	e.SetAttr("id", id)
}

func (e *Element) SetClass(class string) {
	e.SetAttr("class", class)
}

func (e *Element) SetBgColor(color string) {
	e.Style["background-color"] = color
}

func (e *Element) SetBorder(color string) {
	e.Style["border"] = "1.5px solid " + color
	e.Style["border-radius"] = "15px"
}

func (e *Element) SetContentCenter() {
	e.Style["align-items"] = "center"
	e.Style["display"] = "grid"
	e.Style["text-align"] = "center"
}

func (e *Element) AddChild(w ...WebUI) {
	for _, v := range w {
		e.Children = append(e.Children, v)
	}
}

func (e *Element) Render() string {
	styleString := ""
	for k, v := range e.Style {
		styleString += fmt.Sprintf("%s: %s;", k, v)
	}
	if styleString != "" {
		e.Attributes["style"] = styleString
	}
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
		Style:      map[string]string{},
	}
}

func NewElementWithNoEndTag(tag, value string) *Element {
	return &Element{
		Tag: tag, EndTag: "",
		Value:      value,
		Children:   []WebUI{},
		Attributes: map[string]string{},
		Style:      map[string]string{},
	}
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
func NewDiv4C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-4")
	return d
}
func NewDiv5C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-5")
	return d
}
func NewDiv6C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-6")
	return d
}
func NewDiv7C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-7")
	return d
}
func NewDiv8C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-8")
	return d
}
func NewDiv9C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-9")
	return d
}
func NewDiv10C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-10")
	return d
}
func NewDiv11C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-11")
	return d
}
func NewDiv12C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-12")
	return d
}
func NewDivThird(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-third")
	return d
}
func NewDivHalf(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-half")
	return d
}
func NewDivFull(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-full")
	return d
}

func NewCardThird(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-third content-card")
	return d
}
func NewCardHalf(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-half content-card")
	return d
}
func NewCardFull(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-full content-card")
	return d
}

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

// imageType is like png or jpg
func NewImage(imageType string, data []byte) *Element {
	im := NewElementWithNoEndTag("img", "")
	dataStr := base64.StdEncoding.EncodeToString(data)
	im.SetAttr("src", "data:image/"+imageType+";base64,"+dataStr)
	return im
}

func NewImageFromFile(imageType, filename string) *Element {
	im := NewElementWithNoEndTag("img", "")
	data, err := os.ReadFile(filename)
	dataStr := ""
	if err != nil {
		imageType = "png"
		dataStr = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAAXNSR0IArs4c6QAAAARzQklUCAgICHwIZIgAAABzSURBVBhXnZCxCsAgDETPycFBcHT2/7/ETxA3HQUHByclLdIUOoRmOi4vFxIVY1zOOWit8VVzTrTWoHLOi0QIAdbaF9t7R0oJFKRKKcsYcxkcPhB5Y4wb9N6DNyiWD9ZaH5CaBybN0/+BotWiY8TvkT58A7q9hBee+NnzAAAAAElFTkSuQmCC"
	} else {
		dataStr = base64.StdEncoding.EncodeToString(data)
	}
	im.SetAttr("src", "data:image/"+imageType+";base64,"+dataStr)
	return im
}

func NewImageFromLink(imageType, imgLink string) *Element {
	im := NewElementWithNoEndTag("img", "")
	im.SetAttr("src", imgLink)
	return im
}
