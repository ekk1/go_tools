package config

type City interface {
	AddResource(r ResourceType, num int64) bool
	Update()
}

type CityType string

const (
	CityTypeNormal  CityType = "normal"
	CityTypeOutpost CityType = "outpost"
)
