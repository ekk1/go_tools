package config

type Building interface {
	AssignUnit(u UnitType, num int64) bool
	RemoveUnit(u UnitType, num int64) bool
	Update()
}

type BuildingType string

const (
	BuildingTypeFarm    BuildingType = "farm"
	BuildingTypeStorage BuildingType = "storage"
)
