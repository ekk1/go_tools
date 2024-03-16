package webui

import (
	"fmt"
)

var BaseHTML string = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>%s</title>
  <link rel="icon" href="data:;base64,iVBORw0KGgo=">
  <style>
  html {
    font-family: monospace;
    font-size: large;
  }
  body {
    margin: 0;
    width: 100%%;
    height: 100%%;
  }
  * {
    box-sizing: border-box;
  }
  h1 {
    color: DarkMagenta;
  }
  table {
    margin-top: 20px;
  }
  table, th, td {
    border: 1px solid black;
    border-collapse: collapse;
  }
  input {
    margin: 8px;
    font-family: monospace;
    background-color: rgba(255, 250, 240, 0.4);
    border: 1px solid;
    border-radius: 20px;
    padding: 2px 12px 2px;
  }
  th, td {
    padding: 5px;
    border: 1px solid #ddd;
  }
  tr:nth-child(even) {background-color: #f2f2f2;}
  tr:hover {background-color: #ddd;}
  th {
    padding-top: 12px;
    padding-bottom: 12px;
    text-align: left;
    background-color: #04AA6D;
    color: white;
  }
  div {
    padding: 10px;
    margin: 5px;
  }
  .btn {
    color: black;
    margin: 5px;
    font-family: monospace;
    background-color: rgba(0,0,0,0);
    border: 2px solid;
    border-color: rgba(0,0,0,0.5);
    border-radius: 15px;
    padding: 6px 20px;
    text-align: center;
    text-decoration: none;
    display: inline-block;
    font-size: 16px;
    transition: 0.4s;
  }
  .btn:hover {
    background-color: #04AA6D;
    border-color: rgba(0,0,0,0);
    color: white;
  }
  .nav-btn {
    font-weight: bold;
    font-size: 1.12em;
    color: black;
    margin: 1px;
    border: 0;
    border-radius: 0;
    padding: 9px 20px;
    border-bottom: 3px solid rgba(255,20,147,0.2);
  }
  .nav-btn:hover {
    background-color: rgba(255,182,193,0.25);
    color: black;
  }
  .nav-btn-selected {
    background-color: rgba(255,182,193,0.25);
    color: black;
  }
  .nav-btn-div {
    padding: 5px;
    margin-top: 20px;
    display: grid;
    align-items: center;
  }
  .content-card {
    background-color: rgba(220,220,220,0.4);
    border-radius: 30px;
    padding: 10px 30px 25px;
  }
  .ccdiv {
    display: block;
  }
  fieldset {
    border: 1.5px groove rgba(169,169,169,0.3);
    border-radius: 15px;
  }
  [class*="div-"] { grid-column: span 1; }
  @media screen and (min-width: 600px) {
    .ccdiv {
      display: grid;
      grid-template-columns: repeat(6,1fr);
      .div-2 { grid-column: span 2; }
      .div-3 { grid-column: span 3; }
      .div-4 { grid-column: span 4; }
      .div-5 { grid-column: span 5; }
      .div-6 { grid-column: span 6; }
      .div-7 { grid-column: span 6; }
      .div-8 { grid-column: span 6; }
      .div-9 { grid-column: span 6; }
      .div-10 { grid-column: span 6; }
      .div-11 { grid-column: span 6; }
      .div-12 { grid-column: span 6; }
      .div-third { grid-column: span 3; }
      .div-half { grid-column: span 6; }
      .div-full { grid-column: span 6; }
      .div-nav  { grid-column: span 6; }
      .div-cc   { grid-column: span 6; }
    }
  }
  @media screen and (min-width: 768px) {
    .ccdiv {
      grid-template-columns: repeat(8,1fr);
      .div-2 { grid-column: span 2; }
      .div-3 { grid-column: span 3; }
      .div-4 { grid-column: span 3; }
      .div-5 { grid-column: span 4; }
      .div-6 { grid-column: span 4; }
      .div-7 { grid-column: span 5; }
      .div-8 { grid-column: span 6; }
      .div-9 { grid-column: span 6; }
      .div-10 { grid-column: span 7; }
      .div-11 { grid-column: span 7; }
      .div-12 { grid-column: span 8; }
      .div-third { grid-column: span 3; }
      .div-half { grid-column: span 4; }
      .div-full { grid-column: span 8; }
      .div-nav  { grid-column: span 2; }
      .div-cc   { grid-column: span 6; }
    }
  }
  @media screen and (min-width: 1200px) {
    .ccdiv {
      grid-template-columns: repeat(12,1fr);
      .div-2 { grid-column: span 2; }
      .div-3 { grid-column: span 3; }
      .div-4 { grid-column: span 4; }
      .div-5 { grid-column: span 5; }
      .div-6 { grid-column: span 6; }
      .div-7 { grid-column: span 7; }
      .div-8 { grid-column: span 8; }
      .div-9 { grid-column: span 9; }
      .div-10 { grid-column: span 10; }
      .div-11 { grid-column: span 11; }
      .div-12 { grid-column: span 12; }
      .div-third { grid-column: span 4; }
      .div-half { grid-column: span 6; }
      .div-full { grid-column: span 12; }
      .div-nav  { grid-column: span 2; }
      .div-cc   { grid-column: span 10; }
    }
  }
  </style>
</head>
<body>
%s
</body>
</html>`

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
	return fmt.Sprintf(BaseHTML, u.Title, body)
}

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
	nn := NewColumnDiv()
	navDiv := NewDiv(nn)
	navDiv.SetClass("div-nav")
	contentDiv := NewDiv()
	contentDiv.SetClass("div-cc")
	baseDiv := NewColumnDiv(navDiv, contentDiv)
	baseDiv.SetBgColor("rgba(255,248,220,0.2)")
	baseDiv.Style["height"] = "100%"
	baseDiv.Style["width"] = "100%"
	baseDiv.Style["margin"] = "0"
	baseDiv.Style["padding"] = "15px"
	b.Base.AddChild(baseDiv)
	b.NavPane = nn
	b.ContentPane = contentDiv
	navDiv.Style["border-radius"] = "30px"
	navDiv.SetBgColor("rgba(224,255,255,0.90)")
	contentDiv.Style["border-radius"] = "30px"
	contentDiv.SetBgColor("rgba(0,255,127,0.1)")
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
	sectionContainer := NewColumnDiv(ui...)
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
