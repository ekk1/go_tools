package webui

import (
	"fmt"
)

// <div>
type Div struct {
	Children []WebUI
	Tag      string
}

func NewDiv() *Div {
	return &Div{
		Children: []WebUI{},
		Tag:      "<div>",
	}
}

func (u *Div) AddChild(ui ...WebUI) {
	for _, w := range ui {
		u.Children = append(u.Children, w)
	}
}

func (u *Div) SetID(id string) {
	u.Tag = fmt.Sprintf("<div id=\"%s\">", id)
}

func (u *Div) Render() string {
	ret := u.Tag
	for _, w := range u.Children {
		ret += w.Render()
	}
	ret += "</div>"
	return ret
}
