package minikv

import (
	"strconv"
	"testing"
	"time"
)

func TestBenchMiniKV(t *testing.T) {
	kv, err := NewKV("test", 0)
	if err != nil {
		t.Fatal(err)
	}
	runTime := 1 * 1024 * 1024
	dataKey := "key"
	dataValue := "value"
	t1 := time.Now()
	for v := 0; v < runTime; v++ {
		kv.Set(dataKey+"_"+strconv.Itoa(v), dataValue)
	}
	t2 := time.Now()
	tDelta := t2.Sub(t1)
	t.Log("Set (Diff key):", float64(runTime)/tDelta.Seconds(), "IOPS")

	t1 = time.Now()
	for v := 0; v < runTime; v++ {
		kv.Set(dataKey, dataValue)
	}
	t2 = time.Now()
	tDelta = t2.Sub(t1)
	t.Log("Set (Same key):", float64(runTime)/tDelta.Seconds(), "IOPS")

	kv.Save()
}

func TestMiniKV(t *testing.T) {
	kv, err := NewKV("test2", 0)
	if err != nil {
		t.Fatal(err)
	}
	testData := map[string]string{
		"k1": "v1",
		"k2": "v1",
		"k3": "v2",
		"k4": "v2",
		"k5": "v3",
		"k6": "v4",
	}
	for k, v := range testData {
		kv.Set(k, v)
	}

	//kv.Save()
	t.Log(kv.Exists("k1"))
	t.Log(kv.Exists("k2"))
	t.Log(kv.Exists("k8"))
	t.Log(kv.Exists("11"))
	t.Log(kv.Get("k1"))
	t.Log(kv.Get("k2"))

	kv.Delete("k1")
	t.Log(kv.Exists("k1"))
	t.Log(kv.Exists("k2"))

	t.Log(kv.Get("k1"))
	t.Log(kv.Get("k2"))
	t.Log(kv.Get("k8"))
	t.Log(kv.Get("11"))

	kv2, _ := NewKV("test2", 0)
	for k, v := range testData {
		kv2.Set(k, v)
	}
	t.Log(kv2.Get("k1"))
	t.Log(kv2.Get("k2"))
	t.Log(kv2.Get("k8"))
	t.Log(kv2.Get("11"))
	if err := kv2.Load(); err != nil {
		t.Log(err)
	}
	t.Log(kv2.Exists("k1"))
	t.Log(kv2.Exists("k2"))
	t.Log(kv2.Exists("k8"))
	t.Log(kv2.Exists("11"))
	t.Log(kv2.Get("k1"))
	t.Log(kv2.Get("k2"))
	t.Log(kv2.Get("k8"))
	t.Log(kv2.Get("11"))
}
