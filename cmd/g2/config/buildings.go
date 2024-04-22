package config

import "go_tools/cmd/g2/event"

type Building interface {
	AssignUnit(u UnitType, num int64) bool
	RemoveUnit(u UnitType, num int64) bool
	Info() string
	Actions() []string
	Execute(p *event.PlayerEvent) string
	Update()
}

type BuildingType string

const (
	BuildingTypeFarm    BuildingType = "farm"
	BuildingTypeStorage BuildingType = "storage"
)
