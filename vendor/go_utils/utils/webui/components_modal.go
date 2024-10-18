package webui

func NewModal(name, id string, w ...WebUI) *Element {
	modal := NewDiv()
	modal.SetID(id)
	modal.SetClass("w3-modal")
	modal.Style["z-index"] = "6"
	modalContent := NewDiv()
	modalContent.SetClass("w3-modal-content w3-animate-top w3-card-4")
	modalContent.SetClass("w3-round-xlarge w3-sand")

	modalHeader := NewElement("header", "")
	modalHeader.SetClass("w3-container")
	modalHeader.Style["background-color"] = "#04AA6D"
	modalHeader.Style["color"] = "white"
	modalHeader.Style["border-radius"] = "16px 16px 0 0"

	modalCloseBtn := NewElement("span", "&times;")
	modalCloseBtn.SetClass("w3-button w3-display-topright")
	modalCloseBtn.SetAttr("onclick",
		"document.getElementById('"+id+"').style.display='none'",
	)
	modalCloseBtn.Style["border-radius"] = "0 16px 0 0"
	headerText := NewText(name)
	headerText.Style["color"] = "white"
	headerText.Style["font-size"] = "30px"
	headerText.Style["font-weight"] = "bold"
	headerText.Style["margin"] = "10px"
	modalHeader.AddChild(modalCloseBtn, headerText)
	modalContent.AddChild(modalHeader)
	modalContent.AddChild(w...)

	modal.AddChild(modalContent)

	return modal
}
