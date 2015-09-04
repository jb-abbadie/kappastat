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

type ChatEntry struct {
	name   string
	sender string
	time   time.Time
	text   string
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
