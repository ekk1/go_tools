package main

import "go_utils/utils"

type MainObject interface {
	Next()
}

type City struct {
	UUID   string
	CoordX int64
	CoordY int64

	Buildings    map[string]Building
	WorkingUnits map[Unit]int64
	ArmyUnits    map[Unit]int64

	WorkingTeams map[string]Team

	Population    int64
	MaxPopulation int64

	Food        int64
	Storage     map[Resource]int64
	StorageSize int64
	StorageCap  int64
}

func (c *City) AddBuilding(name string, b Building) {
	c.Buildings[name] = b
}

func (c *City) AddTeam(name string) {
	c.WorkingTeams[name] = &WorkingTeam{}
}

func (c *City) AssignWorkingUnitsToBuilding(u Unit, num int64, b string) bool {
	if curNum, ok := c.WorkingUnits[u]; ok && curNum >= num {
		c.WorkingUnits[u] -= num
		return c.Buildings[b].AssignUnit(u, num)
	}
	return false
}

func (c *City) Next() {
	for n, b := range c.Buildings {
		utils.LogPrintDebug(c.UUID, "working on", n)
		b.Next()
	}
	for n, t := range c.WorkingTeams {
		utils.LogPrintDebug(c.UUID, "working on", n)
		t.Next()
	}
}
