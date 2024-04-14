package main

import (
	"fmt"
)

type BlockType uint16

const (
	BlockTypeWater BlockType = iota
	BlockTypeLand
)

type PendingEventType uint16

const (
	PendingEventBuild PendingEventType = iota
	PendingEventPlant
)

var PendingEventChannel = make(chan PendingEventType, 1)
var PendingEventSelection = make(chan byte, 1)

type WorldBlock struct {
	Type     BlockType
	Unit     Unit
	Building Building
	Items    Items
}

func (b *WorldBlock) Show() string {
	if b.Items != nil {
		return b.Items.Show()
	}
	return " . "
}

func (b *WorldBlock) Info() string {
	if b.Items != nil {
		return b.Items.Info()
	}
	ret := ""
	select {
	case e := <-PendingEventChannel:
		switch e {
		case PendingEventPlant:
			ret = "Planting:\n"
			for no, v := range GlobalItemConfigTable.PlantList {
				ret += fmt.Sprintf("\t[%d] %s\n", no, v)
				no++
			}
		default:
			ret = "Unknown pending event"
		}
	default:
		ret = "Empty land\n[p] to plant"
	}
	return ret
}

func (b *WorldBlock) Build(l Building) {
	b.Building = l
}
func (b *WorldBlock) Plant(p PlantType) {
	pp := NewPlant(p)
	b.Items = pp
	GlobalItemList = append(GlobalItemList, pp)
}
