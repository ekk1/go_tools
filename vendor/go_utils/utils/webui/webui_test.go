package webui

import (
	"go_utils/utils"
	"os"
	"testing"
	"time"
)

func TestUI(t *testing.T) {
	t.Log("Testing")

	base := NewNavBase("test")

	base.AddNavItem("Index", "#")
	base.AddNavItem("Config", "#")
	base.AddNavItem("About", "#")
	base.AddNavItem("Links", "#")

	base.CurrentNavItem = "Index"

	hh := NewHeader("test", "h3")
	text := NewText("test1123")
	testPre := NewPreText("tet\n123")

	form1 := NewForm(
		"http://127.0.0.1:5000/post", "test",
		NewTextInputWithValue("test", "value"),
		NewCheckBox("OK"),
		NewRadioInput("radio1", "A"),
		NewRadioInput("radio1", "B"),
		NewRadioInput("radio1", "C"),
		NewTextAreaInput("page"),
		NewSubmitBtn("submit", "submit1"),
	)
	form2 := NewForm(
		"http://127.0.0.1:5000/post", "test",
		NewTextInputWithValue("test2", "value"),
		NewCheckBox("OK2"),
		NewRadioInput("radio2", "A"),
		NewRadioInput("radio2", "B"),
		NewRadioInput("radio2", "C"),
		NewDateInput("date", time.Now(), time.Now().Add(-time.Hour*24*3)),
		NewDateTimeInput("date", 60, time.Now(), time.Now().Add(-time.Hour*24*3)),
		NewDateTimeInput("date", 3600, time.Now(), time.Now().Add(-time.Hour*24*3)),
		NewSubmitBtn("submit", "submit1"),
		NewSubmitBtn("cancel", "submit2"),
	)

	chart := NewChart(200)
	chart.AddData(20, 40)
	chart.AddData(40, 10)
	chart.AddData(60, 90)
	chart.AddData(80, 70)

	paneModal := NewRow(NewCardHalf(hh, text, testPre))
	paneModal.AddChild(NewCardHalf(form1))
	testModal := NewModal("test", "testmodal", paneModal)
	modalBtn := NewBtn("Open")
	modalBtn.SetOpenModal("testmodal")

	pane1 := NewCardHalf(hh, text, testPre, modalBtn)
	pane1.AddChild(testModal)
	pane2 := NewCardHalf(chart)
	pane3 := NewCardHalf(hh, form1)
	pane4 := NewCardHalf(hh, form2)

	pane6 := NewCardThird(hh, form1)

	table := NewTable(
		NewTableRow(true, "Proxy", "Now"),
		NewTableRow(false, "Proxy", "Now"),
		NewTableRow(false, "Proxy", "Now"),
		NewTableRow(false, "Proxy", "Now"),
		NewTableRow(false, "Proxy", "Now"),
		NewTableRow(false, "Proxy", "Now"),
	)
	pane5 := NewCardHalf(table)

	pageTitle := NewCardRest(
		NewHeader("Title", "h2"),
	)
	pageTitle.SetContentCenter()

	paneAbout := NewCardRest(
		NewHeader("About", "h2"),
		NewPreText("this is text\njdwioajwd\n\n\n\n\n\n"),
	)
	paneAbout.SetContentCenter()

	img := NewImageFromLink("placeholder.jpg")
	img.Style["width"] = "7em"
	img.Style["margin-bottom"] = "20px"

	singleUnit := NewDisplayUnit(
		NewHeader("test1", "h4"),
		NewText("test1"),
		img,
		NewBR(),
		NewLinkBtn("test1", "#"),
	)
	simpleUnit := NewDisplayUnit(
		NewHeader("test1", "h4"),
		NewText("test1"),
		NewBR(),
		NewLinkBtn("test1", "#"),
	)
	card := NewCardRest(singleUnit, simpleUnit)
	for i := 0; i < 10; i++ {
		if i%3 == 0 {
			card.AddChild(simpleUnit)
		}
		card.AddChild(singleUnit)
	}
	card.SetFlex()

	base.AddSection("", pageTitle)
	base.AddSection("About", paneAbout)
	base.AddSection("Units", card)
	base.AddSection("section 1", pane1, pane3)
	base.AddSection("section 2", pane2, pane4)
	base.AddSection("section 3", pane5, pane5)
	base.AddSection("section 4", pane6, pane6, pane6)

	utils.LogPrintError(os.MkdirAll("output", 0755))
	utils.LogPrintError(os.WriteFile("output/test.html", []byte(base.Render()), 0644))

	base2 := NewNavBase("Index")
	base2.AddNavItem("Index", "#")
	base2.CurrentNavItem = "Index"
	base2.AddSection("About", paneAbout)
	utils.LogPrintError(os.WriteFile("output/index.html", []byte(base2.Render()), 0644))
}
