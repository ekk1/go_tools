package main

import (
	"fmt"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	m := NewMap(10000, 10000, 1)

	for {
		fmt.Println("Finish Generate: ", m.SizeX)
		time.Sleep(1 * time.Second)
	}
}
