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
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 1fr;
  }
  @media screen and (max-width: 800px) {
    .ccdiv {
      display: grid;
      grid-template-columns: auto;
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
