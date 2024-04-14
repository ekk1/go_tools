package main

import (
	_ "embed"
	"encoding/json"
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

type ItemConfigTable struct {
	ResourceTable map[ResourceType]*ResourceProperty `json:"resources"`
	PlantTable    map[PlantType]*PlantProperty       `json:"plants"`
}
type ResourceProperty struct {
	Size int64 `json:"size"`
}
type PlantProperty struct {
	Outputs map[ResourceType]*PlantOutput `json:"outputs"`
	Max     int64                         `json:"max"`
	Speed   int64                         `json:"speed"`
}
type PlantOutput struct {
	Need   int64 `json:"need"`
	Output int64 `json:"out"`
}

var GlobalItemConfigTable *ItemConfigTable

func LoadConfigTable() error {
	return json.Unmarshal(itemsJson, GlobalItemConfigTable)
}

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
	return fmt.Sprintf("%s: %d", t.Type, t.Grown)
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
