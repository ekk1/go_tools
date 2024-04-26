package config

import (
	_ "embed"
	"encoding/json"
)

type GlobalConfig struct {
	BuildingList []BuildingType                   `json:"build_list"`
	PlantList    []ResourceType                   `json:"plant_list"`
	Resources    map[ResourceType]*ResourceConfig `json:"resource"`
	Units        map[UnitType]*UnitConfig         `json:"unit"`
	Buildings    map[BuildingType]*BuildingConfig `json:"building"`
	Cities       map[CityType]*CityConfig         `json:"city"`
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
	ResourceSize      int64   `json:"size"`
	ResourceValue     int64   `json:"value"`
	ResourceMineSpeed float64 `json:"speed"`
	ResourceOutput    int64   `json:"output"`
}

type UnitConfig struct {
	UnitPopulation     int64   `json:"pop"`
	UnitConsumeFood    int64   `json:"consume"`
	UnitWorkSpeed      float64 `json:"workspeed"`
	UnitMoveSpeed      float64 `json:"movespeed"`
	UnitLoadCapability float64 `json:"load"`
}

type BuildingConfig struct {
	BuildingConstructResources map[ResourceType]int64 `json:"resource"`
	BuildingConstructWork      float64                `json:"work"`
	BuildingMaxWorkingUnits    int64                  `json:"maxunits"`
	BuildingSize               int64                  `json:"size"`
}

type CityConfig struct {
	CityConstructResources map[ResourceType]int64 `json:"resource"`
	CityConstructWork      float64                `json:"work"`
	CityMaxPopulation      int64                  `json:"population_cap"`
	CityStorageCap         int64                  `json:"storage_cap"`
	CityBuildingCap        int64                  `json:"building_cap"`
	CityScanRange          int64                  `json:"scan_range"`
}
