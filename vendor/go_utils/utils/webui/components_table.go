package webui

func NewTable(w ...WebUI) *Element {
	t := NewElement("table", "")
	t.Value = "<tbody>"
	t.ValueEnd = "</tbody>"
	t.AddChild(w...)
	return t
}

func NewTableRow(header bool, data ...string) *Element {
	rowType := "td"
	if header {
		rowType = "th"
	}
	tr := NewElement("tr", "")
	for _, v := range data {
		th := NewElement(rowType, v)
		tr.AddChild(th)
	}
	return tr
}
