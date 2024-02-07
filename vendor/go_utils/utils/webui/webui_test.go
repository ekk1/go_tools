package webui

import (
	"os"
	"testing"
)

func TestUI(t *testing.T) {
	t.Log("Testing")

	base := NewBase("Test")

	cc := NewColumnDiv()
	cc2 := NewColumnDiv()
	cc3 := NewColumnDiv()

	div1 := NewDiv4C()
	div1.SetBgColor("cyan")

	div2 := NewDiv4C()
	div2.SetID("22")
	div2.SetBgColor("pink")

	p1 := NewText("test")

	h1 := NewHeader("test", "h1")
	div1.AddChild(p1, p1, h1)
	div2.AddChild(h1, p1)

	t1 := NewTable()
	t1.AddChild(NewTableRow(true, "header1 test", "header2 test", "h3"))
	t1.AddChild(NewTableRow(false, "h1", "h2", "adwad"))

	div1.AddChild(t1)

	cc.AddChild(div1, div2)

	f1 := NewForm("http://127.0.0.1:5000/post", "test",
		NewTextInput("test"),
	)
	f2 := NewForm("http://127.0.0.1:5000/post", "test2",
		NewTextInput("test88"),
	)

	f1.AddChild(NewTextInput("test91"))
	f1.AddChild(NewTextInputWithValue("test2", "value2"))
	f1.AddChild(NewCheckBox("test3"))
	f1.AddChild(NewRadioInput("test4", "value2"))
	f1.AddChild(NewRadioInput("test4", "value3"))
	f1.AddChild(NewRadioInput("test4", "value4"))
	f1.AddChild(NewRadioInput("test4", "value5"))
	f1.AddChild(NewSubmitBtn("name", "id"))
	f2.AddChild(NewSubmitBtn("name", "id"))

	cc2.AddChild(NewDiv4C(f1))
	cc2.AddChild(NewDiv4C(f2))

	cc3.AddChild(NewDiv6C(NewTable(
		NewTableRow(true, "test1", "test2"),
		NewTableRow(false, "test1", "test2"),
	)))

	base.AddChild(cc, cc2, cc3)
	base.AddChild(NewImageFromFile("jpg", "test.jpg"))

	t.Log(base.Render())
	os.WriteFile("test.html", []byte(base.Render()), 0644)
}
