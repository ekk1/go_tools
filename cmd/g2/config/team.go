package config

import "go_tools/cmd/g2/event"

type Team interface {
	Info() string
	Actions() []string
	Execute(p *event.PlayerEvent) string
	Update()
	AssignUnit(u UnitType, num int64) bool
}

type TeamType string

const (
	TeamTypeWorkerTeam TeamType = "worker_team"
)
