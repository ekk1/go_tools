package city

import (
	"fmt"
	"go_tools/cmd/g2/building"
	"go_tools/cmd/g2/config"
	"go_tools/cmd/g2/event"
	"go_tools/cmd/g2/team"
	"go_utils/utils"
	"strconv"
)

const (
	CityCommandBuild                = "build"
	CityCommandAssignUnitToBuilding = "assign_unit_to_building"
	CityCommandAssignUnitToTeam     = "assign_unit_to_team"
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

func (c *NormalCity) AddBuilding(name string, bType config.BuildingType) bool {
	// c.Buildings[name] = b
	if _, ok := c.Buildings[name]; ok {
		return false
	}
	if c.MaxBuildingSize-c.CurrentBuildingSize < config.Config.Buildings[bType].BuildingSize {
		return false
	}
	switch bType {
	case config.BuildingTypeFarm:
		f := building.NewFarm(c)
		c.Buildings[name] = f
		c.CurrentBuildingSize += config.Config.Buildings[bType].BuildingSize
	}
	return true
}

func (c *NormalCity) AddTeam(name string) bool {
	if _, ok := c.WorkingTeams[name]; ok {
		return false
	}
	c.WorkingTeams[name] = &team.WorkingTeam{}
	return true
}

func (c *NormalCity) AssignWorkingUnitsToBuilding(u config.UnitType, num int64, b string) bool {
	if curNum, ok := c.WorkingUnits[u]; ok && curNum >= num {
		c.WorkingUnits[u] -= num
		return c.Buildings[b].AssignUnit(u, num)
	}
	return false
}

func (c *NormalCity) RemoveWorkingUnitsFromBuilding(u config.UnitType, num int64, b string) bool {
	if curNum, ok := c.WorkingUnits[u]; ok && curNum >= num {
		c.WorkingUnits[u] -= num
		return c.Buildings[b].AssignUnit(u, num)
	}
	return false
}

func (c *NormalCity) AddResource(r config.ResourceType, num int64) bool {
	utils.LogPrintDebug("Adding", r, num)
	pendingResourceSize := config.Config.Resources[r].ResourceSize * num
	if pendingResourceSize+c.CurrentStored > c.MaxStorageCap {
		return false
	}
	c.Storage[r] += num
	c.CurrentStored += pendingResourceSize
	return true
}

func (c *NormalCity) AssignUnit(u config.UnitType, num int64) bool {
	if c.CurrentPopulation+num > c.MaxPopulation {
		return false
	}
	c.CurrentPopulation += num
	c.FoodConsume += config.Config.Units[u].UnitConsumeFood * num
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
	if c.Food < 0 { // City no enough food, stop working
		// Add function for workers to starve ?
		return
	}

	c.Food -= c.FoodConsume
	if c.Food < 0 {
		c.Food = 0
	}

	if c.RemainConstructWorks > 0 { // City not finished building, skip
		for u, num := range c.WorkingUnits {
			c.RemainConstructWorks -= (float64(num) * config.Config.Units[u].UnitWorkSpeed)
			if c.RemainConstructWorks < 0 {
				c.RemainConstructWorks = 0
			}
		}
		return
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
	return []string{CityCommandBuild, CityCommandAssignUnitToBuilding}
}

func (c *NormalCity) Info() string {
	ret := ""
	ret += fmt.Sprintf("CityType: %s	Name: %s\n", config.CityTypeNormal, c.UUID)
	ret += fmt.Sprintf("Coord: x: %d y: %d\n", c.CoordX, c.CoordY)
	if c.RemainConstructWorks != 0 {
		ret += fmt.Sprintf("RemainWork: %f\n", c.RemainConstructWorks)
	}
	ret += fmt.Sprintf("Buildings: 	%d/%d\n", c.CurrentBuildingSize, c.MaxBuildingSize)
	ret += fmt.Sprintf("Population: %d/%d\n", c.CurrentPopulation, c.MaxPopulation)
	ret += fmt.Sprintf("Units: 		%v\n", c.WorkingUnits)
	ret += fmt.Sprintf("Units: 		%v\n", c.Storage)
	ret += fmt.Sprintf("Food: 		%d(-%d)\n", c.Food, c.FoodConsume)
	ret += fmt.Sprintf("Gold: 		%d\n", c.Gold)
	ret += fmt.Sprintf("Stored: 	%d/%d\n", c.CurrentStored, c.MaxStorageCap)
	for n, b := range c.Buildings {
		ret += fmt.Sprintf("Building: %s\n", n)
		ret += fmt.Sprintf("%s\n", b.Info())
	}
	return ret
}

func (c *NormalCity) Execute(e *event.PlayerEvent) string {
	switch e.ActionType {
	case event.PlayerEventTypeCity:
		switch e.Command {
		case CityCommandBuild:
			if _, ok := config.Config.Buildings[config.BuildingType(e.Param2)]; !ok {
				return "Building not supported"
			}
			if !c.AddBuilding(e.Param1, config.BuildingType(e.Param2)) {
				return "No enough space or duplicated building name"
			}
			return "Started building " + e.Param2 + " Name: " + e.Param1
		case CityCommandAssignUnitToBuilding:
			if _, ok := c.WorkingUnits[config.UnitType(e.Param1)]; !ok {
				return "No " + e.Param1 + " availiable"
			}
			if _, ok := c.Buildings[e.Param3]; !ok {
				return "Building " + e.Param3 + " not exists"
			}
			num, err := strconv.ParseInt(e.Param2, 10, 64)
			if err != nil {
				return "Cannot parse " + e.Param2 + " as int"
			}

			if c.AssignWorkingUnitsToBuilding(config.UnitType(e.Param1), num, e.Param3) {
				return "Assigned " + e.Param2 + " " + e.Param1 + " to " + e.Param3
			}
			return "failed"
		case CityCommandAssignUnitToTeam:
			if _, ok := c.WorkingUnits[config.UnitType(e.Param1)]; !ok {
				return "No " + e.Param1 + " availiable"
			}
			if _, ok := c.WorkingTeams[e.Param3]; !ok {
				return "Team " + e.Param3 + " not exists"
			}
			num, err := strconv.ParseInt(e.Param2, 10, 64)
			if err != nil {
				return "Cannot parse " + e.Param2 + " as int"
			}

			if c.WorkingTeams[e.Param3].AssignUnit(config.UnitType(e.Param1), num) {
				return "Assigned " + e.Param2 + " " + e.Param1 + " to " + e.Param3
			}
			return "failed"
		}
	case event.PlayerEventTypeBuilding:
		if b, ok := c.Buildings[e.TargetName]; ok {
			return b.Execute(e)
		}
	}
	return ""
}
