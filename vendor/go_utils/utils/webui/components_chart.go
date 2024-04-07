package webui

import (
	"fmt"
	"math"
)

type Chart struct {
	figure   *Element
	dataList *Element
	dataX    []float64
	dataY    []float64
	width    float64
	height   float64
}

func NewChart(width float64) *Chart {
	fig := NewElement("figure", "")
	fig.SetClass("css-chart")
	fig.Style["--widget-size"] = fmt.Sprintf("%.0fpx;", width)

	figDataList := NewElement("ul", "")
	figDataList.SetClass("line-chart")
	fig.AddChild(figDataList)

	c := &Chart{
		figure:   fig,
		dataList: figDataList,
		width:    width,
	}
	return c
}

func (c *Chart) AddData(dataX, dataY float64) {
	c.dataX = append(c.dataX, dataX)
	c.dataY = append(c.dataY, dataY)
}

func (c *Chart) Render() string {
	if len(c.dataX) != len(c.dataY) {
		return "Failed to render chart: x and y num not equal"
	}
	c.dataList.Children = nil

	for i := 0; i < len(c.dataX); i++ {
		li := NewElement("li", "")
		li.Style["--x"] = fmt.Sprintf("%.2fpx", c.dataX[i])
		li.Style["--y"] = fmt.Sprintf("%.2fpx", c.dataY[i])

		lineLength := float64(0)
		lineAngle := float64(0)

		if i+1 < len(c.dataX) {
			lineLength = math.Sqrt(math.Pow(c.dataX[i+1]-c.dataX[i], 2) + math.Pow(c.dataY[i+1]-c.dataY[i], 2))
			if c.dataY[i+1] > c.dataY[i] {
				lineAngle = -math.Acos((c.dataX[i+1]-c.dataX[i])/lineLength) * (180.0 / math.Pi)
			} else {
				lineAngle = math.Acos((c.dataX[i+1]-c.dataX[i])/lineLength) * (180.0 / math.Pi)
			}
		}

		dataLine := NewElement("div", "")
		dataLine.SetClass("line-segment")
		dataLine.Style["--hypotenuse"] = fmt.Sprintf("%.2f", lineLength)
		dataLine.Style["--angle"] = fmt.Sprintf("%.2f", lineAngle)

		dataPoint := NewElement("div", "")
		dataPoint.SetClass("data-point")
		dataPoint.SetAttr("data-value", li.Style["--y"])

		li.AddChild(dataLine, dataPoint)
		c.dataList.AddChild(li)
	}
	return c.figure.Render()
}
