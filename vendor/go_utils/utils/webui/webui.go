package webui

type WebUI interface {
	Render() string
}

func NewNavBar() *Element {
	d := NewDiv()
	return d
}
