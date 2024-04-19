package main

type Unit string

const (
	UnitTypeWorker Unit = "worker"
)

type UnitConfig struct {
	UnitConsumeFood    map[Unit]int64
	UnitWorkSpeed      map[Unit]float64
	UnitLoadCapability map[Unit]float64
}

var GlobalUnitConfig = &UnitConfig{
	UnitConsumeFood: map[Unit]int64{
		UnitTypeWorker: 1,
	},
	UnitWorkSpeed: map[Unit]float64{
		UnitTypeWorker: 1.0,
	},
	UnitLoadCapability: map[Unit]float64{
		UnitTypeWorker: 10.0,
	},
}

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
