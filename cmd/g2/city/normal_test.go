package city

import (
	"go_tools/cmd/g2/config"
	"testing"
)

func TestNormalCity(t *testing.T) {
	config.InitConfig()

	c := NewCity("test", 0, 0)
	c.RemainConstructWorks = 3
	c.Food = 2

	t.Log(c.AssignUnit(config.UnitTypeWorker, 20))

	for i := 0; i < 6; i++ {
		c.Update()
		t.Log(c.Info())
	}
}
