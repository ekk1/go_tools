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
	CustomJS string
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
	if u.CustomJS != "" {
		return fmt.Sprintf(
			BaseHTML, u.Title, W3StyleStr, StyleStr,
			body,
			fmt.Sprintf("<script>%s</script>", u.CustomJS),
		)
	}
	return fmt.Sprintf(BaseHTML, u.Title, W3StyleStr, StyleStr, "", body)
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
	navOpenScript := "document.getElementById('mySidebar').style.display = 'block';"
	navOpenScript += "document.getElementById('myOverlay').style.display = 'block';"
	navCloseScript := "document.getElementById('mySidebar').style.display = 'none';"
	navCloseScript += "document.getElementById('myOverlay').style.display = 'none';"
	navPane := NewRow()
	navPane.SetClass("w3-pale-blue w3-card w3-sidebar w3-bar-block w3-collapse w3-animate-left")
	navPane.SetClass("w3-margin w3-round-xlarge")
	navPane.SetID("mySidebar")
	//navDiv := NewColQuarter(navPane)
	navDiv := NewDiv(navPane)
	navPane.Style["height"] = "90vh"
	navPane.Style["width"] = "200px"
	navPane.Style["z-index"] = "5"
	navPane.Style["padding"] = "40px 16px 8px 16px"
	navCloseBtn := NewBtn("&times;")
	navCloseBtn.SetClass("w3-display-topright w3-hide-large")
	navCloseBtn.Style["margin"] = "16px"
	navCloseBtn.Style["font-weight"] = "bold"
	navCloseBtn.SetAttr("onclick", navCloseScript)
	navPane.AddChild(navCloseBtn)

	contentPane := NewRow()
	contentPane.Style["min-height"] = "calc(100vh - 48px)"
	contentPane.Style["margin"] = "16px 16px 16px 280px"
	contentPane.SetClass("w3-sand w3-card w3-main w3-padding w3-round-xlarge")
	navOpenBtn := NewBtn("&#9776;")
	navOpenBtn.SetClass("w3-xlarge")
	navOpenBtn.SetClass("w3-hide-large")
	navOpenBtn.Style["margin"] = "16px"
	navOpenBtn.Style["font-weight"] = "bold"
	navOpenBtn.SetAttr("onclick", navOpenScript)
	contentPane.AddChild(navOpenBtn)

	contentDiv := NewColRest(contentPane)

	overlayDiv := NewDiv()
	overlayDiv.SetClass("w3-overlay")
	overlayDiv.SetAttr("onclick", navCloseScript)
	overlayDiv.Style["cursor"] = "pointer"
	overlayDiv.SetID("myOverlay")

	scrollBtn := NewBtn("")
	scrollArrowUp := NewElement("i", "")
	scrollArrowUp.SetClass("arrow up")
	scrollBtn.AddChild(scrollArrowUp)
	scrollBtn.Style["border"] = "3px solid black"
	scrollBtn.SetID("scroll-btn")
	scrollBtn.SetAttr("onclick", "topFunction()")

	b.Base.CustomJS = `let mybutton = document.getElementById("scroll-btn");
window.onscroll = function() {scrollFunction()};
function scrollFunction() {
  if (document.body.scrollTop > 40 || document.documentElement.scrollTop > 40) {
    mybutton.style.display = "block";
  } else {
    mybutton.style.display = "none";
  }
}
function topFunction() {
  document.body.scrollTop = 0;
  document.documentElement.scrollTop = 0;
}`

	b.Base.AddChild(scrollBtn, navDiv, overlayDiv, contentDiv)
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
