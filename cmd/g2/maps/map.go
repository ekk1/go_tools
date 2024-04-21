package maps

import (
	"go_tools/cmd/g2/config"
)

type Object interface {
}

type Map struct {
	SizeX int64
	SizeY int64

	BlocksX map[int64]Object
	BlocksY map[int64]Object
}

func NewMap(sizeX, sizeY int64) *Map {
	worldMap := &Map{
		SizeX:   sizeX,
		SizeY:   sizeY,
		BlocksX: map[int64]Object{},
		BlocksY: map[int64]Object{},
	}

	return worldMap
}

func (m *Map) AddCity(coordX, coordY int64, city *config.City) {
	m.BlocksX[coordX] = city
	m.BlocksY[coordY] = city
}

func (m *Map) AddResource(coordX, coordY int64, resource config.ResourceType) {
	m.BlocksX[coordX] = resource
	m.BlocksY[coordY] = resource
}
