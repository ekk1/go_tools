package main

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
	ResourceSize      map[Resource]int64   `json:"size"`
	ResourceValue     map[Resource]int64   `json:"value"`
	ResourceMineSpeed map[Resource]float64 `json:"speed"`
}

type UnitConfig struct {
	UnitConsumeFood    map[Unit]int64   `json:"consume"`
	UnitWorkSpeed      map[Unit]float64 `json:"workspeed"`
	UnitMoveSpeed      map[Unit]float64 `json:"movespeed"`
	UnitLoadCapability map[Unit]float64 `json:"load"`
}

type BuildingConfig struct {
	BuildingConstructResources map[BuildingName]map[Resource]int64 `json:"resource"`
	BuildingConstructWork      map[BuildingName]float64            `json:"work"`
	BuildingMaxWorkingUnits    map[BuildingName]int64              `json:"maxunits"`
}
