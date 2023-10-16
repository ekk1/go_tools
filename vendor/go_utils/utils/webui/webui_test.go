package webui

import (
	"os"
	"testing"
)

func TestUI(t *testing.T) {
	t.Log("Testing")

	base := NewBase("Test")

	d1 := NewDiv()
	p1 := NewText("test")

	d2 := NewDiv()
	d2.SetID("22")
	h1 := NewHeader("test", 1)
	d1.AddChild(p1, p1, d2)
	d2.AddChild(h1, p1)

	t1 := NewTable()
	t1.AddHeader("h1", "h2")
	t1.AddRow(NewTableRow([]string{"test1", "test2"}))

	d1.AddChild(t1)

	base.AddChild(d1, d2)

	t.Log(base.Render())
	os.WriteFile("test.html", []byte(base.Render()), 0644)
}
