package webui

import (
	_ "embed"
	"fmt"
)

//go:embed base.html
var BaseHTML string

//go:embed style.css
var StyleStr string

//go:embed w3.css
var W3StyleStr string

type WebUI interface {
	Render() string
}

// Base is a simple base for bare bone page
type Base struct {
	Children []WebUI
	Title    string
}

func NewBase(title string) *Base {
	return &Base{
		Title:    title,
		Children: []WebUI{},
	}
}

func (u *Base) AddChild(ui ...WebUI) {
	for _, w := range ui {
		u.Children = append(u.Children, w)
	}
}

func (u *Base) Render() string {
	body := ""
	for _, w := range u.Children {
		body += w.Render()
	}
	return fmt.Sprintf(BaseHTML, u.Title, W3StyleStr, StyleStr, body)
}

// NavBase is a simple base with nav div and content div
type NavBase struct {
	Base           *Base
	NavPane        *Element
	ContentPane    *Element
	NavItems       map[string]*Element
	CurrentNavItem string
}

func NewNavBase(title string) *NavBase {
	b := &NavBase{
		Base: &Base{
			Title:    title,
			Children: []WebUI{},
		},
		NavItems: make(map[string]*Element),
	}
	navPane := NewRow()
	navPane.SetClass("w3-pale-blue")
	navPane.SetBeautifulDiv()
	navDiv := NewColQuarter(navPane)
	navPane.Style["height"] = "100vh"

	contentPane := NewRow()
	contentPane.SetClass("w3-sand")
	contentPane.SetBeautifulDiv()
	contentDiv := NewColThreeQuarter(contentPane)

	baseDiv := NewRow(navDiv, contentDiv)
	b.Base.AddChild(baseDiv)
	b.NavPane = navPane
	b.ContentPane = contentPane
	return b
}

func (n *NavBase) AddNavItem(title, url string) {
	if _, ok := n.NavItems[title]; ok {
		// already added this nav item, skipping
		return
	}
	btn := NewLinkBtn(title, url)
	btn.SetClass("btn nav-btn")
	btnDiv := NewDiv(btn)
	btnDiv.SetClass("div-full nav-btn-div")
	n.NavItems[title] = btn
	n.NavPane.AddChild(btnDiv)
}

func (n *NavBase) AddSection(title string, ui ...WebUI) {
	sectionContainer := NewRow(ui...)
	if title != "" {
		sectionTitle := NewHeader(title, "h4")
		sectionTitle.Style["margin"] = "1px"
		sectionTitle.Style["padding"] = "1px"
		sectionHeader := NewDiv(sectionTitle)
		sectionHeader.Style["margin"] = "1px"
		sectionHeader.Style["padding"] = "10px 1px 1px 20px"
		n.AddContent(sectionHeader, sectionContainer)
	} else {
		n.AddContent(sectionContainer)
	}
}

func (n *NavBase) AddContent(w ...WebUI) {
	n.ContentPane.AddChild(w...)
}

func (n *NavBase) Render() string {
	ret := ""
	if item, ok := n.NavItems[n.CurrentNavItem]; ok {
		item.SetClass("btn nav-btn nav-btn-selected")
		ret = n.Base.Render()
		item.SetClass("btn nav-btn")
	} else {
		ret = n.Base.Render()
	}
	return ret
}
