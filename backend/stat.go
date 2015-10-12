package main

import (
	"github.com/grsakea/kappastat/common"
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func loopStat(ch chan Message, db *mgo.Database) {
	followed := []string{}

	c := cron.New()

	c.AddFunc("0 * * * * *", func() { computeStat(db, followed, 01*time.Minute) })
	c.AddFunc("0 */5 * * * *", func() { computeStat(db, followed, 05*time.Minute) })
	c.AddFunc("0 */15 * * * *", func() { computeStat(db, followed, 15*time.Minute) })
	c.AddFunc("@hourly", func() { computeStat(db, followed, time.Hour) })
	c.AddFunc("0 0 */12 * * *", func() { computeStat(db, followed, 12*time.Hour) })
	c.AddFunc("@daily", func() { computeStat(db, followed, 24*time.Hour) })

	c.Start()
	for {
		select {
		case msg := <-ch:
			if msg.s == EndBroadcast {
				processBroadcast(db, msg.v)
			} else {
				followed = followedHandler(followed, msg)
			}
		}

	}
}

func computeStat(db *mgo.Database, channels []string, duration time.Duration) {
	to := time.Now()
	from := to.Add(-duration)

	for _, channel := range channels {
		data, err := fetchStatData(db, channel, from, to)
		if err == nil {
			se := processStatData(from, to, duration, channel, data)
			storeStatEntry(db.C("stat_entries"), se)
		}
	}
}

func processBroadcast(db *mgo.Database, channel string) {
	var v kappastat.ViewerCount
	db.C("stat_entries").Find(bson.M{"channel": channel}).Sort("-Time").One(&v)
	log.Print(v)
}

func processStatData(from time.Time, to time.Time, duration time.Duration, channel string, data statData) (ret kappastat.StatEntry) {
	ret.Channel = channel
	ret.Duration = duration
	ret.Start = from
	ret.End = to

	var resultC kappastat.ChatEntry
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

	var result kappastat.ViewerCount
	ret.Viewer = 0
	nbZero := 0
	for data.itV.Next(&result) {
		ret.Viewer += result.Viewer
		if result.Viewer == 0 {
			nbZero--
		}
	}
	ret.Viewer /= data.lenV
	if nbZero > data.lenV {
		ret.NonZeroViewer /= (data.lenV - nbZero)
	} else {
		ret.NonZeroViewer = 0
	}
	return
}
