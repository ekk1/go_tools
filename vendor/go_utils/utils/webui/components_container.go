package webui

func NewDiv(w ...WebUI) *Element {
	d := NewElement("div", "")
	d.AddChild(w...)
	return d
}

func NewRow(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("w3-row")
	return d
}

func NewCol(size string, w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass(size)
	return d
}
func NewColThird(w ...WebUI) *Element {
	return NewCol("w3-third", w...)
}
func NewColTwoThird(w ...WebUI) *Element {
	return NewCol("w3-twothird", w...)
}
func NewColQuarter(w ...WebUI) *Element {
	return NewCol("w3-quarter", w...)
}
func NewColThreeQuarter(w ...WebUI) *Element {
	return NewCol("w3-threequarter", w...)
}
func NewColHalf(w ...WebUI) *Element {
	return NewCol("w3-half", w...)
}
func NewColRest(w ...WebUI) *Element {
	return NewCol("w3-rest", w...)
}

func NewCard(size string, w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("w3-card")
	d.SetBeautifulDiv()
	return NewCol(size, d)
}
func NewCardThird(w ...WebUI) *Element {
	return NewCard("w3-third", w...)
}
func NewCardTwoThird(w ...WebUI) *Element {
	return NewCard("w3-twothird", w...)
}
func NewCardQuater(w ...WebUI) *Element {
	return NewCard("w3-quarter", w...)
}
func NewCardThreeQuater(w ...WebUI) *Element {
	return NewCard("w3-threequarter", w...)
}
func NewCardHalf(w ...WebUI) *Element {
	return NewCard("w3-half", w...)
}
func NewCardRest(w ...WebUI) *Element {
	return NewCard("w3-rest", w...)
}
