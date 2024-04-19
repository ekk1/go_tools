package main

type Building interface {
	AssignUnit(u Unit, num int64) bool
	RemoveUnit(u Unit, num int64) bool
	Next()
}

type FarmCorpType string

const (
	FarmCorpTypeCorn FarmCorpType = "corn"
)

// TODO: Finish farm
type Farm struct {
	ParentCity *City

	Planting       FarmCorpType
	MaxGrown       float64
	GrowSpeed      float64
	CurrentGrown   float64
	ExpectedOutput float64

	UnitNum  int64
	MaxUnits int64
}

func (f *Farm) Next() {
	f.CurrentGrown += f.GrowSpeed
	if f.CurrentGrown > f.MaxGrown {
		f.CurrentGrown -= f.MaxGrown
		// Output to parent city's storage
	}
}

func (f *Farm) AssignUnit(u Unit, num int64) bool {
	if f.UnitNum+num > f.MaxUnits {
		return false
	}
	f.UnitNum += num
	f.GrowSpeed += GlobalUnitConfig.UnitWorkSpeed[u] * float64(num)

	return true
}
