package main

import (
	"time"
)

type Histo struct {
	viewers  map[time.Time]int
	streamer string
}

type ViewerCount struct {
	Channel string
	Time    time.Time
	Viewer  int
}

type ChatEntry struct {
	Channel string
	Sender  string
	Time    time.Time
	Text    string
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
