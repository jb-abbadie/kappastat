package main

import (
	"time"
)

type Histo struct {
	viewers  map[time.Time]int
	streamer string
}

type StreamState struct {
	name   string
	time   time.Time
	viewer int
}

type Signal int

const (
	AddStream Signal = iota
	RemoveStream
	Stop
	Restart
)

type Message struct {
	s Signal
	v string
}
