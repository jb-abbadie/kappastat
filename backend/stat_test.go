package main

import (
	"github.com/grsakea/kappastat/common"
	"gopkg.in/mgo.v2"
	"testing"
	"time"
)

func TestStatProcessing(t *testing.T) {
	session, err := mgo.Dial("localhost")

	if err != nil {
		t.Error("Error connecting to the db")
		t.FailNow()
	}
	db := session.DB("twitch_test")
	db.DropDatabase()

	tim := time.Now()
	ce1 := kappastat.ChatEntry{"test", "test_user", tim, "This is a test message"}
	db.C("chat_entries").Insert(ce1)
	vc1 := kappastat.ViewerCount{"test", tim, 42}
	db.C("viewer_count").Insert(vc1)

	computeStat(db, []string{"test"}, time.Hour)

	var ret kappastat.StatEntry
	db.C("stat_entries").Find(nil).One(&ret)

	if ret.Duration != time.Hour {
		t.Error("Duration incorrect")
	}
	if ret.Messages != 1 {
		t.Error("Number of messages incorrect")
	}
	if ret.Newsub != 0 {
		t.Error("Newsub count incorrect")
	}
	if ret.Resub != 0 {
		t.Error("Resub count incorrect")
	}
	if ret.UniqueChat != 1 {
		t.Error("Unique chat incorrect")
	}
	if ret.Viewer != 42 {
		t.Error("ViewerCount incorrect")
	}
}
