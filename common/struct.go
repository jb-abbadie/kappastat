package kappastat

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
