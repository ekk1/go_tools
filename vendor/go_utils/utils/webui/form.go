package webui

import (
	"fmt"
)

// <form>
type Form struct {
	URL      string
	Method   string
	EncType  string
	Legend   string
	Tag      string
	Children []WebUI
}

func NewForm(u, name string) *Form {
	return &Form{
		URL:     u,
		Method:  "post",
		EncType: "multipart/form-data",
		Legend:  name,
		Tag: fmt.Sprintf(
			"<form action=\"%s\" method=\"%s\" enctype=\"%s\">",
			u, "post",
			"multipart/form-data",
		),
		Children: []WebUI{},
	}
}

func (u *Form) AddChild(ui ...WebUI) {
	for _, w := range ui {
		u.Children = append(u.Children, w)
	}
}

func (u *Form) SetID(id string) {
	u.Tag = fmt.Sprintf(
		"<form action=\"%s\" method=\"%s\" enctype=\"%s\" id=\"%s\">",
		u, "post",
		"multipart/form-data",
		id,
	)
}

func (u *Form) Render() string {
	ret := u.Tag + "<fieldset>"
	ret += fmt.Sprintf("<legend>%s</legend>", u.Legend)
	for _, w := range u.Children {
		ret += w.Render()
	}
	ret += "</fieldset></form>"
	return ret
}

// <input>
type Input struct {
	Text string
	Tag  string
	Type string
	Name string
}

func NewInput(name, inputType, content, id string) *Input {
	return &Input{
		Text: content,
		Tag: fmt.Sprintf(
			"<input type=\"%s\" name=\"%s\" value=\"%s\" id=\"%s\">",
			inputType, name, content, id,
		),
		Type: inputType,
	}
}

func (u *Input) SetID(id string) {
	u.Tag = fmt.Sprintf(
		"<input type=\"%s\" name=\"%s\" value=\"%s\" id=\"%s\">",
		u.Type, u.Name, u.Text, id,
	)
}

func (u *Input) Render() string {
	return u.Tag
}

// <label>
type Label struct {
	Text   string
	Target string
	Tag    string
}

func NewLabel(content string, target string) *Label {
	return &Label{
		Text:   content,
		Tag:    fmt.Sprintf("<label for=\"%s\">", target),
		Target: target,
	}
}

func (u *Label) SetID(id string) {
	u.Tag = fmt.Sprintf("<label id=\"%s\", for=\"%s\">", id, u.Target)
}

func (u *Label) Render() string {
	return fmt.Sprintf("%s%s</label>", u.Tag, u.Text)
}

type GroupTextInput struct {
	Name        string
	Label       *Label
	Input       *Input
	PlaceHolder string
}

func NewGroupTextInput(name string) *GroupTextInput {
	l := NewLabel(name, name)
	i := NewInput(name, "text", "", name)
	return &GroupTextInput{Name: name, Label: l, Input: i, PlaceHolder: ""}
}

func NewGroupTextInputWithPlaceHolder(name string, placeholder string) *GroupTextInput {
	l := NewLabel(name, name)
	i := NewInput(name, "text", placeholder, name)
	return &GroupTextInput{Name: name, Label: l, Input: i, PlaceHolder: placeholder}
}

func (u *GroupTextInput) Render() string {
	return u.Label.Render() + u.Input.Render()
}
