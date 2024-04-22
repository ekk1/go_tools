package player

import "go_tools/cmd/g2/config"

type PlayerStruct struct {
	UUID     string
	CityList map[string]config.City
}
