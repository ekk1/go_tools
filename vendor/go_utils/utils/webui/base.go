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
  <style>
  h1 {
    color: DarkMagenta;
  }
  table, th, td {
    border: 1px solid black;
  }
  input {
    margin: 5px;
  }
  th, td {
    padding: 5px;
  }
  div {
    margin: 10px;
    margin-top: 20px;
  }
  button {
    margin: 5px;
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
