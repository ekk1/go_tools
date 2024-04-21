package team

import (
	"go_tools/cmd/g2/config"
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
}

func (t *WorkingTeam) Update() {

}
