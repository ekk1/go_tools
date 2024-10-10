package webui

import (
	"fmt"
	"slices"
)

type Element struct {
	Tag        string
	Attributes map[string]string
	Value      string
	ValueEnd   string
	EndTag     string
	Style      map[string]string
	Class      []string
	Children   []WebUI
}

func (e *Element) SetAttr(k, v string) {
	e.Attributes[k] = v
}

func (e *Element) SetID(id string) {
	e.SetAttr("id", id)
}

func (e *Element) SetClass(class string) {
	//e.SetAttr("class", class)
	e.Class = append(e.Class, class)
}

func (e *Element) RemoveClass(class string) {
	classIndex := slices.Index(e.Class, class)
	fmt.Println(classIndex)
	if classIndex != -1 {
		e.Class = slices.Delete(e.Class, classIndex, classIndex+1)
	}
}

func (e *Element) SetBgColor(color string) {
	e.Style["background-color"] = color
}

func (e *Element) SetRound() {
	e.SetClass("w3-round-xlarge")
}
func (e *Element) SetPadding() {
	e.SetClass("w3-padding")
}
func (e *Element) SetMargin() {
	e.SetClass("w3-margin")
}
func (e *Element) SetMarginPadding() {
	e.SetClass("w3-margin w3-padding")
}
func (e *Element) SetBeautifulDiv() {
	e.SetClass("w3-margin w3-padding w3-round-xlarge")
}

func (e *Element) SetOpenModal(id string) {
	e.SetAttr("onclick",
		"document.getElementById('"+id+"').style.display='block'",
	)
}

func (e *Element) SetBorder(color string) {
	e.Style["border"] = "1.5px solid " + color
	e.Style["border-radius"] = "15px"
}

func (e *Element) SetContentCenter() {
	e.SetClass("w3-center")
}
func (e *Element) SetContentLeft() {
	e.SetClass("w3-left")
}
func (e *Element) SetContentRight() {
	e.SetClass("w3-right")
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
	classStr := ""
	for _, v := range e.Class {
		classStr += fmt.Sprintf("%s ", v)
	}
	if classStr != "" {
		e.Attributes["class"] = classStr
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
		Class:      []string{},
	}
}

func NewElementWithNoEndTag(tag, value string) *Element {
	return &Element{
		Tag: tag, EndTag: "",
		Value:      value,
		Children:   []WebUI{},
		Attributes: map[string]string{},
		Style:      map[string]string{},
		Class:      []string{},
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
