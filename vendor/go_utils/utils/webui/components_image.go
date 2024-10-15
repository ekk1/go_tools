package webui

import (
	"encoding/base64"
	"os"
)

// imageType is like png or jpg
func NewImage(imageType string, data []byte) *Element {
	im := NewElementWithNoEndTag("img", "")
	dataStr := base64.StdEncoding.EncodeToString(data)
	im.SetAttr("src", "data:image/"+imageType+";base64,"+dataStr)
	return im
}

func NewImageFromFile(imageType, filename string) *Element {
	im := NewElementWithNoEndTag("img", "")
	data, err := os.ReadFile(filename)
	dataStr := ""
	if err != nil {
		imageType = "png"
		dataStr = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAAXNSR0IArs4c6QAAAARzQklUCAgICHwIZIgAAABzSURBVBhXnZCxCsAgDETPycFBcHT2/7/ETxA3HQUHByclLdIUOoRmOi4vFxIVY1zOOWit8VVzTrTWoHLOi0QIAdbaF9t7R0oJFKRKKcsYcxkcPhB5Y4wb9N6DNyiWD9ZaH5CaBybN0/+BotWiY8TvkT58A7q9hBee+NnzAAAAAElFTkSuQmCC"
	} else {
		dataStr = base64.StdEncoding.EncodeToString(data)
	}
	im.SetAttr("src", "data:image/"+imageType+";base64,"+dataStr)
	return im
}

func NewImageFromLink(imgLink string) *Element {
	im := NewElementWithNoEndTag("img", "")
	im.SetAttr("src", imgLink)
	return im
}
