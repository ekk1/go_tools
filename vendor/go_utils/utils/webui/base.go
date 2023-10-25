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
  * {
    box-sizing: border-box;
  }
  h1 {
    color: DarkMagenta;
  }
  table, th, td {
    border: 1px solid black;
    border-collapse: collapse;
  }
  input {
    margin: 5px;
    font-family: monospace;
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
  }
  .btn {
    margin: 5px;
    font-family: monospace;
    background-color: rgba(0,0,0,0);
    border-color: #04AA6D;
    color: black;
    border: 2px solid;
    padding: 8px 16px;
    text-align: center;
    text-decoration: none;
    display: inline-block;
    font-size: 16px;
    border-radius: 8px;
  }
  .btn:hover {
    background-color: #04AA6D;
    color: white;
  }
  .ccdiv {
    display: block;
  }
  [class*="div-"] { grid-column: span 1; }
  @media screen and (min-width: 600px) {
    .ccdiv {
      display: grid;
      grid-template-columns: repeat(6, 1fr);
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
    }
  }
  @media screen and (min-width: 768px) {
    .ccdiv {
      grid-template-columns: repeat(8, 1fr);
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
    }
  }
  @media screen and (min-width: 1200px) {
    .ccdiv {
      grid-template-columns: repeat(12, 1fr);
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
