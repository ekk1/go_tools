package config

import (
	_ "embed"
	"encoding/json"
)

type GlobalConfig struct {
	Resources *ResourceConfig `json:"resource"`
	Units     *UnitConfig     `json:"unit"`
	Buildings *BuildingConfig `json:"building"`
}

var Config *GlobalConfig

//go:embed config.json
var configBytes []byte

func InitConfig() {
	Config = &GlobalConfig{}
	if err := json.Unmarshal(configBytes, Config); err != nil {
		panic(err)
	}
}

type ResourceConfig struct {
	ResourceSize      map[ResourceType]int64   `json:"size"`
	ResourceValue     map[ResourceType]int64   `json:"value"`
	ResourceMineSpeed map[ResourceType]float64 `json:"speed"`
	ResourceOutput    map[ResourceType]int64   `json:"output"`
}

type UnitConfig struct {
	UnitPopulation     map[UnitType]int64   `json:"pop"`
	UnitConsumeFood    map[UnitType]int64   `json:"consume"`
	UnitWorkSpeed      map[UnitType]float64 `json:"workspeed"`
	UnitMoveSpeed      map[UnitType]float64 `json:"movespeed"`
	UnitLoadCapability map[UnitType]float64 `json:"load"`
}

type BuildingConfig struct {
	BuildingConstructResources map[BuildingType]map[ResourceType]int64 `json:"resource"`
	BuildingConstructWork      map[BuildingType]float64                `json:"work"`
	BuildingMaxWorkingUnits    map[BuildingType]int64                  `json:"maxunits"`
	BuildingSize               map[BuildingType]int64                  `json:"size"`
}

type CityConfig struct {
	CityConstructResources map[CityType]map[ResourceType]int64 `json:"resource"`
	CityConstructWork      map[CityType]float64                `json:"work"`
	CityMaxPopulation      map[CityType]int64                  `json:"population"`
	CityStorageCap         map[CityType]int64                  `json:"cap"`
}
