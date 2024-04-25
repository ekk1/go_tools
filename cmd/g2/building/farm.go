package building

import (
	"fmt"
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
		MaxUnits:             config.Config.Buildings[config.BuildingTypeFarm].BuildingMaxWorkingUnits,
		RemainConstructWorks: config.Config.Buildings[config.BuildingTypeFarm].BuildingConstructWork,
		UnitsList:            map[config.UnitType]int64{},
	}
}

func (f *Farm) Update() {
	if f.RemainConstructWorks > 0 {
		// f.RemainConstructWorks -= (float64(f.UnitNum) * )
		for u, num := range f.UnitsList {
			f.RemainConstructWorks -= (float64(num) * config.Config.Units[u].UnitWorkSpeed)
			if f.RemainConstructWorks < 0 {
				f.RemainConstructWorks = 0
			}
		}
		return
	}
	if f.Planting == "" {
		return
	}
	if f.CurrentGrown > f.MaxGrown {
		if f.ParentCity.AddResource(f.Planting, f.ExpectedOutput) {
			f.CurrentGrown -= f.MaxGrown
		}
	} else {
		f.CurrentGrown += f.GrowSpeed * config.Config.Resources[f.Planting].ResourceMineSpeed
	}
}

func (f *Farm) Plant(r config.ResourceType) {
	f.Planting = r
	f.ExpectedOutput = config.Config.Resources[r].ResourceOutput
	f.CurrentGrown = 0
	f.MaxGrown = float64(config.Config.Resources[r].ResourceValue)
}

func (f *Farm) AssignUnit(u config.UnitType, num int64) bool {
	if f.UnitNum+num > f.MaxUnits {
		return false
	}
	f.UnitNum += num
	f.GrowSpeed += config.Config.Units[u].UnitWorkSpeed * float64(num)
	if _, ok := f.UnitsList[u]; !ok {
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
	if e.ActionType != event.PlayerEventTypeBuilding {
		return "Error: not supported action type"
	}
	switch e.Command {
	case "plant":
		if _, ok := config.Config.Resources[config.ResourceType(e.Param1)]; !ok {
			return "Building not supported"
		}
		f.Plant(config.ResourceType(e.Param1))
		return "Planted " + e.Param1
	}
	return ""
}

func (f *Farm) Info() string {
	ret := ""
	ret += fmt.Sprintf("Planting: %s Grown/Max/Speed: %f/%f/%f\n",
		f.Planting, f.CurrentGrown, f.MaxGrown, f.GrowSpeed,
	)
	if f.RemainConstructWorks != 0 {
		ret += fmt.Sprintf("RemainWork:	%f\n", f.RemainConstructWorks)
	}
	ret += fmt.Sprintf("Worker/Max: %d/%d\n", f.UnitNum, f.MaxUnits)
	return ret
}
