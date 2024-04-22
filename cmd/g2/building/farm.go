package building

import (
	"go_tools/cmd/g2/config"
	"go_tools/cmd/g2/event"
)

type Farm struct {
	ParentCity config.City

	Planting       config.ResourceType
	MaxGrown       float64
	GrowSpeed      float64
	CurrentGrown   float64
	ExpectedOutput int64

	RemainConstructWorks float64
	UnitNum              int64
	UnitsList            map[config.UnitType]int64
	MaxUnits             int64
}

func NewFarm(c config.City) *Farm {
	return &Farm{
		ParentCity:           c,
		MaxUnits:             config.Config.Buildings.BuildingMaxWorkingUnits[config.BuildingTypeFarm],
		RemainConstructWorks: config.Config.Buildings.BuildingConstructWork[config.BuildingTypeFarm],
		UnitsList:            map[config.UnitType]int64{},
	}
}

func (f *Farm) Update() {
	if f.RemainConstructWorks > 0 {
		// f.RemainConstructWorks -= (float64(f.UnitNum) * )
		for u, num := range f.UnitsList {
			f.RemainConstructWorks -= (float64(num) * config.Config.Units.UnitWorkSpeed[u])
		}
		return
	}
	if f.CurrentGrown > f.MaxGrown {
		if f.ParentCity.AddResource(f.Planting, f.ExpectedOutput) {
			f.CurrentGrown -= f.MaxGrown
		}
	} else {
		f.CurrentGrown += f.GrowSpeed
	}
}

func (f *Farm) Plant(r config.ResourceType) {
	f.Planting = r
	f.ExpectedOutput = config.Config.Resources.ResourceOutput[r]
	f.CurrentGrown = 0
}

func (f *Farm) AssignUnit(u config.UnitType, num int64) bool {
	if f.UnitNum+num > f.MaxUnits {
		return false
	}
	f.UnitNum += num
	f.GrowSpeed += config.Config.Units.UnitWorkSpeed[u] * float64(num)
	if num, ok := f.UnitsList[u]; !ok {
		f.UnitsList[u] = num
	} else {
		f.UnitsList[u] += num
	}

	return true
}

func (f *Farm) RemoveUnit(u config.UnitType, num int64) bool {
	return true
}

func (f *Farm) Actions() []string {
	return []string{
		"Plant",
	}
}

func (f *Farm) Execute(e *event.PlayerEvent) string {
	return ""
}

func (f *Farm) Info() string {
	return ""
}
