package config

type Team interface {
	Update()
}

type TeamType string

const (
	TeamTypeWorkerTeam TeamType = "worker_team"
)
