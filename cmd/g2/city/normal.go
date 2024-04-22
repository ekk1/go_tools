package city

import (
	"fmt"
	"go_tools/cmd/g2/building"
	"go_tools/cmd/g2/config"
	"go_tools/cmd/g2/event"
	"go_tools/cmd/g2/team"
	"go_utils/utils"
)

type NormalCity struct {
	UUID   string
	CoordX int64
	CoordY int64

	RemainConstructWorks float64

	Buildings           map[string]config.Building
	CurrentBuildingSize int64
	MaxBuildingSize     int64
	WorkingUnits        map[config.UnitType]int64
	ArmyUnits           map[config.UnitType]int64

	WorkingTeams map[string]config.Team

	CurrentPopulation int64
	MaxPopulation     int64
	Food              int64
	FoodConsume       int64

	Gold int64

	Storage       map[config.ResourceType]int64
	CurrentStored int64
	MaxStorageCap int64

	// string is "coordX,coordY"
	NearbyResources map[string]config.ResourceType
	MaxScanRange    int64
}

func NewCity(name string, x, y int64) *NormalCity {
	return &NormalCity{
		UUID:                 name,
		CoordX:               x,
		CoordY:               y,
		RemainConstructWorks: config.Config.Cities[config.CityTypeNormal].CityConstructWork,
		Buildings:            make(map[string]config.Building),
		MaxBuildingSize:      config.Config.Cities[config.CityTypeNormal].CityBuildingCap,
		WorkingUnits:         map[config.UnitType]int64{},
		ArmyUnits:            map[config.UnitType]int64{},
		WorkingTeams:         map[string]config.Team{},
		MaxPopulation:        config.Config.Cities[config.CityTypeNormal].CityMaxPopulation,
		Storage:              map[config.ResourceType]int64{},
		MaxStorageCap:        config.Config.Cities[config.CityTypeNormal].CityStorageCap,
		MaxScanRange:         config.Config.Cities[config.CityTypeNormal].CityScanRange,
		NearbyResources:      map[string]config.ResourceType{},
	}
}

func NewCityNoWork(name string, x, y int64) *NormalCity {
	c := NewCity(name, x, y)
	c.RemainConstructWorks = 0
	return c
}

func (c *NormalCity) AddBuilding(name string, bType config.BuildingType) {
	// c.Buildings[name] = b
	switch bType {
	case config.BuildingTypeFarm:
		f := building.NewFarm(c)
		c.Buildings[name] = f
		// c.CurrentBuildingSize +=
	}
}

func (c *NormalCity) AddTeam(name string) {
	c.WorkingTeams[name] = &team.WorkingTeam{}
}

func (c *NormalCity) AssignWorkingUnitsToBuilding(u config.UnitType, num int64, b string) bool {
	if curNum, ok := c.WorkingUnits[u]; ok && curNum >= num {
		c.WorkingUnits[u] -= num
		return c.Buildings[b].AssignUnit(u, num)
	}
	return false
}

func (c *NormalCity) RemoveWorkingUnitsToBuilding(u config.UnitType, num int64, b string) bool {
	if curNum, ok := c.WorkingUnits[u]; ok && curNum >= num {
		c.WorkingUnits[u] -= num
		return c.Buildings[b].AssignUnit(u, num)
	}
	return false
}

func (c *NormalCity) AddResource(r config.ResourceType, num int64) bool {
	pendingResourceSize := config.Config.Resources[r].ResourceSize * num
	if pendingResourceSize+c.CurrentStored > c.MaxStorageCap {
		return false
	}
	if _, ok := c.Storage[r]; ok {
		c.Storage[r] += num
		c.CurrentStored += pendingResourceSize
	}
	return true
}

func (c *NormalCity) AssignUnit(u config.UnitType, num int64) bool {
	if c.CurrentPopulation+num > c.MaxPopulation {
		return false
	}
	c.CurrentPopulation += num
	c.FoodConsume += config.Config.Units[u].UnitConsumeFood
	if _, ok := c.WorkingUnits[u]; !ok {
		c.WorkingUnits[u] = num
	} else {
		c.WorkingUnits[u] += num
	}
	return true
}

func (c *NormalCity) ScanResources() {

}

func (c *NormalCity) Update() {
	if c.RemainConstructWorks > 0 {
		for u, num := range c.WorkingUnits {
			c.RemainConstructWorks -= (float64(num) * config.Config.Units[u].UnitWorkSpeed)
		}
		return
	}

	if c.Food < 0 {
		return
	}

	c.Food -= c.FoodConsume
	if c.Food < 0 {
		c.Food = 0
	}

	for n, b := range c.Buildings {
		utils.LogPrintDebug(c.UUID, "working on", n)
		b.Update()
	}
	for n, t := range c.WorkingTeams {
		utils.LogPrintDebug(c.UUID, "working on", n)
		t.Update()
	}
}

func (c *NormalCity) Actions() []string {
	return []string{}
}

func (c *NormalCity) Info() string {
	ret := ""
	ret += fmt.Sprintf("CityType:	%s\n", config.CityTypeNormal)
	ret += fmt.Sprintf("UUID: 		%s\n", c.UUID)
	ret += fmt.Sprintf("CoordX: 	%d\n", c.CoordX)
	ret += fmt.Sprintf("CoordY: 	%d\n", c.CoordY)
	ret += fmt.Sprintf("RemainWork: %f\n", c.RemainConstructWorks)
	ret += fmt.Sprintf("Buildings: 	%d\n", c.CurrentBuildingSize)
	ret += fmt.Sprintf("MaxBuilds: 	%d\n", c.MaxBuildingSize)
	ret += fmt.Sprintf("Population: %d\n", c.CurrentPopulation)
	ret += fmt.Sprintf("MaxPop: 	%d\n", c.MaxPopulation)
	ret += fmt.Sprintf("Units: 		%v\n", c.WorkingUnits)
	ret += fmt.Sprintf("Food: 		%d\n", c.Food)
	ret += fmt.Sprintf("FoodRate: 	%d\n", c.FoodConsume)
	ret += fmt.Sprintf("Gold: 		%d\n", c.Gold)
	ret += fmt.Sprintf("Stored: 	%d\n", c.CurrentStored)
	ret += fmt.Sprintf("MaxStore: 	%d\n", c.MaxStorageCap)
	return ret
}

func (c *NormalCity) Execute(e *event.PlayerEvent) string {
	return ""
}
