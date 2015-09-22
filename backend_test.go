package main

import (
	"github.com/grsakea/kappastat/backend"
	"github.com/mrshankly/go-twitch/twitch"
	"gopkg.in/mgo.v2"
	"net/http"
	"testing"
	"time"
)

func Test_backend(t *testing.T) {
	client := twitch.NewClient(&http.Client{})
	tmp, _ := mgo.Dial("127.0.0.1")
	db := tmp.DB("twitch_test")
	db.DropDatabase()
	coll := db.C("follow")
	opt := &twitch.ListOptions{
		Limit:  5,
		Offset: 0,
	}
	l, _ := client.Streams.List(opt)
	for _, item := range l.Streams {
		tmp, err := client.Users.User(item.Channel.Name)
		if err != nil {
			t.Fatal(err)
		}
		err = coll.Insert(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}

	expected, err := coll.Find(nil).Count()

	if err != nil {
		t.Fatal(err)
	}

	c := backend.SetupController("twitch_test")
	go c.Loop()
	time.Sleep(30 * time.Second)
	if len(c.ListStreams()) != expected {
		t.Fatal("Failed loading tracked list from the DB")
	}
}
