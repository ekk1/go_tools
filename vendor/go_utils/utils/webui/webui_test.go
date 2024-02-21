package webui

import (
	"os"
	"testing"
)

func TestUI(t *testing.T) {
	t.Log("Testing")

	base := NewNavBase("test")

	base.AddNavItem("Index", "#")
	base.AddNavItem("Config", "#")
	base.AddNavItem("About", "#")
	base.AddNavItem("Links", "#")

	base.CurrentNavItem = "Index"

	//hh := NewHeader("test", "h3")
	//text := NewText("test1123")
	//testPre := NewPreText("tet\n123")

	form1 := NewForm("http://127.0.0.1:5000/post", "test")
	form1.AddChild(
		NewTextInputWithValue("test", "value"),
		NewCheckBox("OK"),
		NewRadioInput("radio1", "A"),
		NewRadioInput("radio1", "B"),
		NewRadioInput("radio1", "C"),
		NewSubmitBtn("submit", "submit1"),
	)
	form2 := NewForm("http://127.0.0.1:5000/post", "test")
	form2.AddChild(
		NewTextInputWithValue("test2", "value"),
		NewCheckBox("OK2"),
		NewRadioInput("radio2", "A"),
		NewRadioInput("radio2", "B"),
		NewRadioInput("radio2", "C"),
		NewSubmitBtn("submit", "submit1"),
		NewSubmitBtn("cancel", "submit2"),
	)

	//pane1 := NewCardHalf(hh, text, testPre)
	//pane2 := NewCardHalf(hh, text, testPre)
	//pane3 := NewCardHalf(hh, form1)
	//pane4 := NewCardHalf(hh, form2)
	//base.AddSection("section 1", pane1, pane3)
	//base.AddSection("section 2", pane2, pane4)

	paneAbout := NewCardFull(
		NewHeader("About", "h2"),
		NewPreText("this is text\njdwioajwd\n\n\n\n\n\n\n\n\n\n\n\n\n\n"),
	)
	paneAbout.SetContentCenter()
	base.AddContent(paneAbout)

	os.WriteFile("output/test.html", []byte(base.Render()), 0644)
}
