package webui

import (
	"fmt"
)

// <h[1-6]>
type Header struct {
	Level int
	Text  string
	Tag   string
}

func NewHeader(content string, level int) *Header {
	actualLevel := 0
	switch {
	case level < 1:
		actualLevel = 1
	case level > 6:
		actualLevel = 6
	default:
		actualLevel = level
	}
	return &Header{
		Text:  content,
		Level: actualLevel,
		Tag:   fmt.Sprintf("<h%d>", actualLevel),
	}
}
func (u *Header) SetID(id string) {
	u.Tag = fmt.Sprintf("<h%d id=\"%s\">", u.Level, id)
}

func (u *Header) Render() string {
	return fmt.Sprintf("%s%s</h%d>", u.Tag, u.Text, u.Level)
}

// <p>
type Text struct {
	Text string
	Tag  string
}

func NewText(content string) *Text {
	return &Text{
		Text: content,
		Tag:  "<p>",
	}
}
func (u *Text) SetID(id string) {
	u.Tag = fmt.Sprintf("<p id=\"%s\">", id)
}

func (u *Text) Render() string {
	return fmt.Sprintf("%s%s</p>", u.Tag, u.Text)
}

// <a>
// Link is a <a>
type Link struct {
	Text string
	URL  string
	Tag  string
}

func NewLink(content string, link string) *Link {
	return &Link{
		Text: content,
		URL:  link,
		Tag:  fmt.Sprintf("<a href=\"%s\">", link),
	}
}

func (u *Link) SetID(id string) {
	u.Tag = fmt.Sprintf("<a href=\"%s\" id=\"%s\">", u.URL, id)
}

func (u *Link) Render() string {
	return fmt.Sprintf("%s%s</a>", u.Tag, u.Text)
}

// <button>
// Button is a <button> with <a> inside, currently this is the only one and required
type Button struct {
	ID   string
	Tag  string
	Link *Link
}

func NewButton(l *Link) *Button {
	return &Button{
		Tag:  "<button>",
		Link: l,
	}
}

func (u *Button) SetID(id string) {
	u.Tag = fmt.Sprintf("<button id=\"%s\">", id)
}

func (u *Button) Render() string {
	return fmt.Sprintf("%s%s</button>", u.Tag, u.Link.Render())
}

// <br>
type BR struct {
}

func NewBR() *BR {
	return &BR{}
}

func (u *BR) Render() string {
	return "<br>"
}
