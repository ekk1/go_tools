package webui

import (
	"fmt"
	"testing"
)

func TestClass(t *testing.T) {
	e := NewElement("div", "")
	e.SetClass("test2")
	e.SetClass("test<:F4>")
	e.SetClass("test3")
	e.SetClass("test")

	e.RemoveClass("test")
	e.RemoveClass("test")
	e.RemoveClass("test")
	e.RemoveClass("test")

	fmt.Println(e.Class)

}
