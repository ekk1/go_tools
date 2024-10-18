package webui

func NewLoginPage(loginURL, title string) *Base {
	lBase := NewBase("Login")

	contentPane := NewRow()
	contentPane.Style["min-height"] = "calc(100vh - 48px)"
	contentPane.Style["margin"] = "16px"
	contentPane.SetClass("w3-sand w3-card w3-padding w3-round-xlarge")

	formPane := NewForm(loginURL, title,
		NewTextInput("Username"),
		NewPasswordInput("Password"),
		NewSubmitBtn("Login", "login"),
	)

	placeHolder := NewCardThird()
	placeHolder.e.Style["visibility"] = "hidden"

	loginCard := NewCardThird(NewHeader(title, "h2"), formPane)
	loginCard.e.Style["margin-top"] = "calc(10vh + 20px)"
	loginCard.e.Style["margin-bottom"] = "calc(10vh + 20px)"

	contentPane.AddChild(placeHolder, loginCard)
	contentPane.SetContentCenter()

	lBase.AddChild(contentPane)
	return lBase
}
