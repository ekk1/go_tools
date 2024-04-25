package city

import (
	"go_tools/cmd/g2/config"
	"go_tools/cmd/g2/event"
	"go_utils/utils"
	"testing"
)

func TestNormalCity(t *testing.T) {
	config.InitConfig()
	utils.SetLogLevelByVerboseFlag(1)
	t.Log(config.Config.Resources[config.ResourceTypeCorn])

	c := NewCity("test", 0, 0)
	c.RemainConstructWorks = 6
	c.Food = 100

	t.Log(c.AssignUnit(config.UnitTypeWorker, 5))

	for i := 0; i < 3; i++ {
		c.Update()
		t.Log(c.Info())
	}

	e1 := &event.PlayerEvent{
		ActionType: event.PlayerEventTypeCity,
		CityName:   "test",
		TargetName: "test",
		Command:    "build",
		Param1:     "f1",
		Param2:     "farm",
	}
	t.Log(c.Execute(e1))

	e2 := &event.PlayerEvent{
		ActionType: event.PlayerEventTypeCity,
		CityName:   "test",
		TargetName: "test",
		Command:    "assign_unit_to_building",
		Param1:     "worker",
		Param2:     "5",
		Param3:     "f1",
	}
	t.Log(c.Execute(e2))

	for i := 0; i < 3; i++ {
		c.Update()
		t.Log(c.Info())
	}

	e3 := &event.PlayerEvent{
		ActionType: event.PlayerEventTypeBuilding,
		CityName:   "test",
		TargetName: "f1",
		Command:    "plant",
		Param1:     "corn",
	}
	t.Log(c.Execute(e3))

	for i := 0; i < 3; i++ {
		c.Update()
		t.Log(c.Info())
	}

}
