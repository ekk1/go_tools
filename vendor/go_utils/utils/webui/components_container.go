package webui

func NewDiv(w ...WebUI) *Element {
	d := NewElement("div", "")
	d.AddChild(w...)
	return d
}

func NewColumnDiv(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("ccdiv")
	return d
}

func NewDiv2C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-2")
	return d
}
func NewDiv3C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-3")
	return d
}
func NewDiv4C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-4")
	return d
}
func NewDiv5C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-5")
	return d
}
func NewDiv6C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-6")
	return d
}
func NewDiv7C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-7")
	return d
}
func NewDiv8C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-8")
	return d
}
func NewDiv9C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-9")
	return d
}
func NewDiv10C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-10")
	return d
}
func NewDiv11C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-11")
	return d
}
func NewDiv12C(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-12")
	return d
}
func NewDivThird(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-third")
	return d
}
func NewDivHalf(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-half")
	return d
}
func NewDivFull(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-full")
	return d
}

func NewCardThird(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-third content-card")
	return d
}
func NewCardHalf(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-half content-card")
	return d
}
func NewCardFull(w ...WebUI) *Element {
	d := NewDiv(w...)
	d.SetClass("div-full content-card")
	return d
}
