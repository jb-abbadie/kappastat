package backend

import (
	"github.com/mrshankly/go-twitch/twitch"
	"gopkg.in/mgo.v2"
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

type Controller struct {
	config      Config
	infosChat   chan ChatEntry
	infosViewer chan ViewerCount
	cViewer     chan Message
	cChat       chan Message
	cStat       chan Message
	tracked     map[string]bool
	storage     StorageController
	twitchAPI   *twitch.Client
}

type StorageController struct {
	db     *mgo.Database
	views  *mgo.Collection
	chat   *mgo.Collection
	follow *mgo.Collection
}

type StatEntry struct {
	Channel string

	Duration time.Duration
	Start    time.Time
	End      time.Time

	Resub         int
	Newsub        int
	Messages      int
	UniqueChat    int
	MostUsedTerm  []string
	Viewer        int
	NonZeroViewer int
}
