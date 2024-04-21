package city

import (
	"go_tools/cmd/g2/building"
	"go_tools/cmd/g2/config"
	"go_tools/cmd/g2/team"
	"go_utils/utils"
)

type NormalCity struct {
	UUID   string
	CoordX int64
	CoordY int64

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
	pendingResourceSize := config.Config.Resources.ResourceSize[r] * num
	if pendingResourceSize+c.CurrentStored > c.MaxStorageCap {
		return false
	}
	if _, ok := c.Storage[r]; ok {
		c.Storage[r] += num
		c.CurrentStored += pendingResourceSize
	}
	return true
}

func (c *NormalCity) Update() {
	for n, b := range c.Buildings {
		utils.LogPrintDebug(c.UUID, "working on", n)
		b.Update()
	}
	for n, t := range c.WorkingTeams {
		utils.LogPrintDebug(c.UUID, "working on", n)
		t.Update()
	}
}
