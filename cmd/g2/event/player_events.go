package event

type PlayerEvent struct {
	ActionType string `json:"action"`
	TargetName string `json:"name"`
	Command    string `json:"command"`
	Param1     string `json:"param1"`
	Param2     string `json:"param2"`
	Param3     string `json:"param3"`
	Param4     string `json:"param4"`
	Finished   chan string
	Output     string
}
