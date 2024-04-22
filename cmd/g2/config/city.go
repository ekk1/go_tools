package config

import "go_tools/cmd/g2/event"

type City interface {
	AssignUnit(u UnitType, num int64) bool
	AddResource(r ResourceType, num int64) bool
	Info() string
	Actions() []string
	Execute(p *event.PlayerEvent) string
	Update()
}

type CityType string

const (
	CityTypeNormal  CityType = "normal"
	CityTypeOutpost CityType = "outpost"
)
