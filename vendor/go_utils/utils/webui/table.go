package webui

import (
	"fmt"
)

// <table>
type Table struct {
	Headers []string
	Rows    []*TableRow
	Tag     string
}

func NewTable() *Table {
	return &Table{
		Headers: []string{},
		Rows:    []*TableRow{},
		Tag:     "<table>",
	}
}

func (u *Table) SetID(id string) {
	u.Tag = fmt.Sprintf("<table id=\"%s\">", id)
}

func (u *Table) AddRow(r ...*TableRow) {
	for _, w := range r {
		u.Rows = append(u.Rows, w)
	}
}

func (u *Table) Clear() {
	u.Rows = []*TableRow{}
}

func (u *Table) ClearHeader() {
	u.Headers = []string{}
}

func (u *Table) AddHeader(headers ...string) {
	for _, h := range headers {
		u.Headers = append(u.Headers, h)
	}
}

func (u *Table) Render() string {
	ret := u.Tag + "<tr>"
	for _, v := range u.Headers {
		ret += "<th>" + v + "</th>"
	}
	ret += "</tr>"
	for _, v := range u.Rows {
		ret += v.Render()
	}
	ret += "</table>"
	return ret
}

type TableRow struct {
	Data []string
}

func NewTableRow(data []string) *TableRow {
	return &TableRow{Data: data}
}

func (u *TableRow) AddData(data ...string) {
	for _, d := range data {
		u.Data = append(u.Data, d)
	}
}

func (u *TableRow) Clear() {
	u.Data = []string{}
}

func (u *TableRow) Render() string {
	ret := "<tr>"
	for _, v := range u.Data {
		ret += "<td>" + v + "</td>"
	}
	ret += "</tr>"
	return ret
}
