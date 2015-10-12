package kappastat

import (
	"time"
)

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

type Broadcast struct {
	Start time.Time
	End   time.Time

	Channel string
	Game    string
	Name    string

	AverageViewership int
	MinViewership     int
	MaxViewership     int

	Sub   int
	ReSub int
}
