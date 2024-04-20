package main

type Building interface {
	AssignUnit(u Unit, num int64) bool
	RemoveUnit(u Unit, num int64) bool
	Next()
}

type BuildingName string

const (
	BuildingNameFarm BuildingName = "farm"
)

// TODO: Finish farm
type Farm struct {
	ParentCity *City

	Planting       Resource
	MaxGrown       float64
	GrowSpeed      float64
	CurrentGrown   float64
	ExpectedOutput float64

	RemainConstructWorks float64
	UnitNum              int64
	UnitsList            map[Unit]int64
	MaxUnits             int64
}

func NewFarm(c *City) *Farm {
	return &Farm{
		ParentCity:           c,
		MaxUnits:             Config.Buildings.BuildingMaxWorkingUnits[BuildingNameFarm],
		RemainConstructWorks: Config.Buildings.BuildingConstructWork[BuildingNameFarm],
		UnitsList:            map[Unit]int64{},
	}
}

func (f *Farm) Next() {
	if f.RemainConstructWorks > 0 {
		// f.RemainConstructWorks -= (float64(f.UnitNum) * )
		for u, num := range f.UnitsList {
			f.RemainConstructWorks -= (float64(num) * Config.Units.UnitWorkSpeed[u])
		}
		return
	}
	f.CurrentGrown += f.GrowSpeed
	if f.CurrentGrown > f.MaxGrown {
		f.CurrentGrown -= f.MaxGrown
		// Output to parent city's storage
	}
}

func (f *Farm) AssignUnit(u Unit, num int64) bool {
	if f.UnitNum+num > f.MaxUnits {
		return false
	}
	f.UnitNum += num
	f.GrowSpeed += Config.Units.UnitWorkSpeed[u] * float64(num)

	return true
}
