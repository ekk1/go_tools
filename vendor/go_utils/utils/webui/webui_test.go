package webui

import (
	"os"
	"testing"
)

func TestUI(t *testing.T) {
	t.Log("Testing")

	base := NewBase("Test")

	div1 := NewDiv()
	p1 := NewText("test")

	div2 := NewDiv()
	div2.SetID("22")

	h1 := NewHeader("test", "h1")
	div1.AddChild(p1, p1, h1)
	div2.AddChild(h1, p1)

	t1 := NewTable()
	t1.AddChild(NewTableRow(true, "h1", "h2"))
	t1.AddChild(NewTableRow(false, "h1", "h2"))

	div1.AddChild(t1)

	base.AddChild(div1, div2)

	f1 := NewForm("/", "test",
		NewTextInput("test"),
	)
	f2 := NewForm("/", "test2",
		NewTextInput("test88"),
	)

	f1.AddChild(NewTextInput("test91"))
	f1.AddChild(NewTextInputWithValue("test2", "value2"))
	f1.AddChild(NewCheckBox("test3"))
	f1.AddChild(NewRadioInput("test4", "value2"))
	f1.AddChild(NewRadioInput("test4", "value3"))
	f1.AddChild(NewRadioInput("test4", "value4"))
	f1.AddChild(NewRadioInput("test4", "value5"))

	base.AddChild(NewDiv(f1))
	base.AddChild(NewDiv(f2))
	base.AddChild(NewDiv(NewTable(
		NewTableRow(true, "test1", "test2"),
		NewTableRow(false, "test1", "test2"),
	)))

	t.Log(base.Render())
	os.WriteFile("test.html", []byte(base.Render()), 0644)
}
