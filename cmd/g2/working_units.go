package main

type Unit string

const (
	UnitTypeWorker Unit = "worker"
)

type Team interface {
	Next()
}

type WorkingTeam struct {
	TargetCoordX   int64
	TargetCoordY   int64
	TargetResource Resource

	TargetDistance  float64
	TeamMoveSpeed   float64
	TeamCurrentMove float64
	TeamWorkSpeed   float64

	TeamLoadCap     int64
	TeamFoodConsume int64
	TeamFoodNum     int64
}

func (t *WorkingTeam) Next() {

}
