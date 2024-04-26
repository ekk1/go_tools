package team

import (
	"go_tools/cmd/g2/config"
	"go_tools/cmd/g2/event"
)

type WorkingTeam struct {
	TargetCoordX   int64
	TargetCoordY   int64
	TargetResource config.ResourceType

	TargetDistance  float64
	TeamMoveSpeed   float64
	TeamCurrentMove float64
	TeamWorkSpeed   float64

	TeamLoadCap     int64
	TeamFoodConsume int64
	TeamFoodNum     int64

	TeamMaxUnitNum int64
	TeamUnitNum    int64
	TeamUnits      map[config.UnitType]int64
}

func (t *WorkingTeam) Update() {

}

func (t *WorkingTeam) Actions() []string {
	return []string{}
}

func (t *WorkingTeam) Execute(p *event.PlayerEvent) string {
	return "OK"
}

func (t *WorkingTeam) Info() string {
	return ""
}

func (t *WorkingTeam) AssignUnit(u config.UnitType, num int64) bool {
	if t.TeamUnitNum+num > t.TeamMaxUnitNum {
		return false
	}
	t.TeamUnitNum += num
	t.TeamFoodConsume += config.Config.Units[u].UnitConsumeFood * num

	if _, ok := t.TeamUnits[u]; !ok {
		t.TeamUnits[u] = num
	} else {
		t.TeamUnits[u] += num
	}
	return true
}
