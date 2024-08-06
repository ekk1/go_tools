package main

import "sync"

type GameState int64

const (
	GameStateWaiting GameState = iota
	GameStateCalling
	GameStatePlaying
	GameStateDealing
)

type Game struct {
	mu              sync.Mutex
	players         []*Player
	gameState       GameState
	currentTurn     int
	landlord        int
	lastPlayedCards []string // 添加这个字段
	bottomCards     []string
}
