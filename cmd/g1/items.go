package main

import (
	_ "embed"
	"fmt"
)

//go:embed items_prop.json
var itemsJson []byte

type (
	ResourceType string
	PlantType    string
)

const (
	ResourceTypeApple  ResourceType = "apple"
	ResourceTypeWood   ResourceType = "wood"
	ResourceTypeCotton ResourceType = "cotton"
	ResourceTypeRice   ResourceType = "rice"
)
const (
	PlantTypeAppleTree PlantType = "apple_tree"
	PlantTypeWoodTree  PlantType = "wood_tree"
	PlantTypeCotton    PlantType = "cotton"
	PlantTypeRice      PlantType = "rice"
)

type Resource struct {
	Type ResourceType
	Size int64
	Num  int64
}

type Plant struct {
	Type      PlantType
	Grown     int64
	MaxGrown  int64
	GrowSpeed int64
	Outputs   map[ResourceType]*PlantOutput
}

func NewPlant(p PlantType) *Plant {
	return &Plant{
		Type:      p,
		Grown:     0,
		MaxGrown:  GlobalItemConfigTable.PlantTable[p].Max,
		GrowSpeed: GlobalItemConfigTable.PlantTable[p].Speed,
		Outputs:   GlobalItemConfigTable.PlantTable[p].Outputs,
	}
}

func (t *Plant) Show() string {
	return fmt.Sprintf(" T%d", t.Grown/10)
}

func (t *Plant) Next() {
	t.Grow()
}

func (t *Plant) Info() string {
	return fmt.Sprintf(
		"Type: %s\nGrown: %d / %d\n[h]: harvest, [x]: delete\n",
		t.Type, t.Grown, t.MaxGrown,
	)
}

func (t *Plant) Grow() {
	if t.Grown < t.MaxGrown {
		t.Grown += t.GrowSpeed
	}
	if t.Grown > t.MaxGrown {
		t.Grown = t.MaxGrown
	}
}

func (t *Plant) Harvest(r ResourceType) *Resource {
	if out, ok := t.Outputs[r]; !ok {
		return nil
	} else {
		if t.Grown < out.Need {
			return nil
		}
		t.Grown -= out.Need
		return &Resource{
			Type: r,
			Num:  out.Output,
		}
	}
}
