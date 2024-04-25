package event

type PlayerEventType string

const (
	PlayerEventTypeCity     PlayerEventType = "city"
	PlayerEventTypeBuilding PlayerEventType = "building"
)

type PlayerEvent struct {
	ActionType PlayerEventType `json:"action"`
	CityName   string          `json:"city_name"`
	TargetName string          `json:"target_name"` // This should be the city name if PlayerEventTypeCity
	Command    string          `json:"command"`     // This is the action
	Param1     string          `json:"param1"`
	Param2     string          `json:"param2"`
	Param3     string          `json:"param3"`
	Param4     string          `json:"param4"`
	Finished   chan string
	Output     string
}
