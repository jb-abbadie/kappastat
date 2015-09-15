package backend

import (
	"gopkg.in/mgo.v2"
	"log"
	"strings"
	"time"
)

type statData struct {
	itC  *mgo.Iter
	lenC int
	itV  *mgo.Iter
	lenV int
}

func loopStat(c chan Message, db *mgo.Database) {
	followed := []string{}
	oneMinute := time.NewTicker(time.Minute).C
	tenMinute := time.NewTicker(10 * time.Minute).C
	oneHour := time.NewTicker(time.Hour).C
	oneDay := time.NewTicker(24 * time.Hour).C

	for {
		select {
		case msg := <-c:
			followed = followedHandler(followed, msg)
		case <-oneMinute:
			go computeStat(db, followed, time.Minute)
		case <-tenMinute:
			go computeStat(db, followed, 10*time.Minute)
		case <-oneHour:
			go computeStat(db, followed, time.Hour)
		case <-oneDay:
			go computeStat(db, followed, 24*time.Hour)
		}

	}
}

func computeStat(db *mgo.Database, channels []string, duration time.Duration) {
	to := time.Now()
	from := to.Add(-duration)

	for _, channel := range channels {
		data := fetchStatData(db, channel, from, to)
		se := processStatData(from, to, duration, channel, data)
		storeStatEntry(db.C("stat_entries"), se)
	}
}

func processStatData(from time.Time, to time.Time, duration time.Duration, channel string, data statData) (ret StatEntry) {
	ret.Channel = channel
	ret.Duration = duration
	ret.Start = from
	ret.End = to

	var resultC ChatEntry
	uniqueChatter := make(map[string]bool)
	termUsed := make(map[string]int)

	for data.itC.Next(&resultC) {
		if resultC.Sender == "twitchnotify" {
			if strings.Contains(resultC.Text, "just") {
				ret.Newsub += 1
			} else if strings.Contains(resultC.Text, "months") {
				ret.Resub += 1
			}
		} else {
			ret.Messages += 1

			for _, i := range strings.Split(resultC.Text, " ") {
				termUsed[i] += 1
			}

			_, present := uniqueChatter[resultC.Sender]
			if !present {
				uniqueChatter[resultC.Sender] = true
			}
		}
	}
	ret.UniqueChat = len(uniqueChatter)

	var result ViewerCount
	ret.Viewer = 0
	for data.itV.Next(&result) {
		ret.Viewer += result.Viewer
	}
	ret.Viewer /= data.lenV
	return
}
