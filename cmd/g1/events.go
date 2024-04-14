package main

type GameEvent uint8

const (
	EventCameraUP GameEvent = iota
	EventCameraDown
	EventCameraLeft
	EventCameraRight
	EventGameStop
	EventPlant
	EventDoPlant
)
