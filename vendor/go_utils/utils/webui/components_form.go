package webui

import (
	"fmt"
	"go_utils/utils"
)

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
