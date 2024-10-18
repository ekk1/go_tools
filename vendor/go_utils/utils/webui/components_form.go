package webui

import (
	"go_utils/utils"
	"strconv"
	"time"
)

func NewForm(targetURL, name string, e ...WebUI) *Element {
	f := NewElement("form", "")
	//f.Value = fmt.Sprintf("<fieldset><legend>%s</legend>", name)
	//f.ValueEnd = "</fieldset>"
	f.SetAttr("action", targetURL)
	f.SetAttr("method", "post")
	f.SetAttr("enctype", "multipart/form-data")
	f.AddChild(e...)
	return f
}

func NewInput(name, inputType, content, id string) *Element {
	i := NewElementWithNoEndTag("input", "")
	i.SetAttr("type", inputType)
	i.SetAttr("name", name)
	i.SetAttr("value", content)
	i.SetAttr("id", id)
	i.SetAttr("autocomplete", "off")
	return i
}

func NewSubmitBtn(name string, id string) *GroupElement {
	s := NewInput(name, "submit", name, id)
	s.SetClass("btn")
	return NewGroupElement(s, NewBR())
}

func NewLabel(content string, target string) *Element {
	l := NewElement("label", content)
	l.SetAttr("for", target)
	return l
}

func NewGroupInput(name, value, iType string) *GroupElement {
	idSuffix := utils.RandomString(5)
	txInput := NewInput(name, iType, value, name+"-"+idSuffix)
	txInput.SetClass("text-input")
	return NewGroupElement(
		NewLabel(name, name+"-"+idSuffix),
		txInput, NewBR(),
	)
}

func NewTextInputWithValue(name, value string) *GroupElement {
	return NewGroupInput(name, value, "text")
}
func NewTextInput(name string) *GroupElement {
	return NewTextInputWithValue(name, "")
}

func NewPasswordInput(name string) *GroupElement {
	return NewGroupInput(name, "", "password")
}

func NewTextAreaInput(name string) *GroupElement {
	idSuffix := utils.RandomString(5)
	area := NewElement("textarea", "")
	area.SetAttr("name", name)
	area.SetAttr("id", name+"-"+idSuffix)
	return NewGroupElement(
		NewLabel(name, name+"-"+idSuffix), NewBR(),
		area, NewBR(),
	)
}

func NewCheckBox(name string) *GroupElement {
	idSuffix := utils.RandomString(5)
	return NewGroupElement(
		NewLabel(name, name+"-"+idSuffix),
		NewInput(name, "checkbox", name, name+"-"+idSuffix), NewBR(),
	)
}

func NewRadioInput(name, value string) *GroupElement {
	return NewGroupElement(
		NewInput(name, "radio", value, name+"-"+value),
		NewLabel(value, name+"-"+value), NewBR(),
	)
}

func NewDateInput(name string, maxDate, minData time.Time) *GroupElement {
	curTime := time.Now()
	curTimeStr := curTime.Format(time.DateOnly)
	idSuffix := utils.RandomString(5)
	dateInput := NewInput(name, "date", curTimeStr, name+"-"+idSuffix)
	dateInput.SetAttr("max", maxDate.Format(time.DateOnly))
	dateInput.SetAttr("min", minData.Format(time.DateOnly))
	return NewGroupElement(
		NewLabel(name, name+"-"+idSuffix),
		dateInput, NewBR(),
	)
}

// Step is second, default to 60
func NewDateTimeInput(name string, step int64, maxDate, minData time.Time) *GroupElement {
	curTime := time.Now()
	curTimeStr := curTime.Format("2006-01-02T15:04:05")
	idSuffix := utils.RandomString(5)
	dateInput := NewInput(name, "datetime-local", curTimeStr, name+"-"+idSuffix)
	dateInput.SetAttr("max", maxDate.Format("2006-01-02T15:04:05"))
	dateInput.SetAttr("min", minData.Format("2006-01-02T15:04:05"))
	dateInput.SetAttr("step", strconv.FormatInt(step, 10))
	return NewGroupElement(
		NewLabel(name, name+"-"+idSuffix),
		dateInput, NewBR(),
	)
}
