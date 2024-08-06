package main

import (
	"net"
	"sync"
)

type PlayerStatus int64

const (
	PlayerStatusWaiting PlayerStatus = iota
	PlayerStatusPlaying
	PlayerStatusDisconnected
	PlayerStatusRejected
)

type Player struct {
	conn   net.Conn
	name   string
	hand   []string
	status PlayerStatus
	mu     sync.Mutex
}
