package main

import "encoding/json"

type Unit interface {
	Next()
	Show() string
	Info() string
}

type Building interface {
	Next()
	Show() string
	Info() string
}

type Items interface {
	Next()
	Show() string
	Info() string
}

type ItemConfigTable struct {
	ResourceTable map[ResourceType]*ResourceProperty `json:"resources"`
	PlantTable    map[PlantType]*PlantProperty       `json:"plants"`

	ResourceList []ResourceType
	PlantList    []PlantType
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
	GlobalItemConfigTable = &ItemConfigTable{}
	err := json.Unmarshal(itemsJson, GlobalItemConfigTable)
	if err != nil {
		return err
	}
	for k := range GlobalItemConfigTable.ResourceTable {
		GlobalItemConfigTable.ResourceList = append(GlobalItemConfigTable.ResourceList, k)
	}
	for k := range GlobalItemConfigTable.PlantTable {
		GlobalItemConfigTable.PlantList = append(GlobalItemConfigTable.PlantList, k)
	}
	return nil
}
