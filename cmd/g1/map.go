package main

import "fmt"

type Map struct {
	SizeX int64
	SizeY int64
	SizeZ int64

	Blocks [][][]*WorldBlock
}

func NewMap(sizeX, sizeY, sizeZ int64) *Map {
	worldMap := &Map{
		SizeX:  sizeX,
		SizeY:  sizeY,
		SizeZ:  sizeZ,
		Blocks: [][][]*WorldBlock{},
	}

	for z := int64(0); z < sizeZ; z++ {
		mapPane := [][]*WorldBlock{}
		for y := int64(0); y < sizeY; y++ {
			mapLine := []*WorldBlock{}
			for x := int64(0); x < sizeX; x++ {
				block := &WorldBlock{
					Type: BlockTypeLand,
				}
				mapLine = append(mapLine, block)
			}
			mapPane = append(mapPane, mapLine)
		}
		worldMap.Blocks = append(worldMap.Blocks, mapPane)
	}
	return worldMap
}

func (m *Map) Render(centerX, centerY int64, viewRange int64) string {
	ret := ""
	for z := int64(0); z < m.SizeZ; z++ {
		// ret += fmt.Sprintf("Layer %d:\n", z)
		mapPane := ""
		startLane := centerY - viewRange
		if startLane < 0 {
			startLane = 0
		}
		endLane := centerY + viewRange
		if endLane > m.SizeY {
			endLane = m.SizeY
		}
		for y := startLane; y < endLane; y++ {
			mapLine := ""
			startPoint := centerX - viewRange
			if startPoint < 0 {
				startPoint = 0
			}
			endPoint := centerX + viewRange
			if endPoint > m.SizeX {
				endPoint = m.SizeX
			}
			for x := startPoint; x < endPoint; x++ {
				if x == centerX && y == centerY {
					mapLine += " + "
					continue
				}
				mapLine += fmt.Sprintf("%v", m.Blocks[z][y][x].Show())
			}
			mapPane += mapLine + "\n\n"
		}
		ret += mapPane + "\n"
	}
	return ret
}
