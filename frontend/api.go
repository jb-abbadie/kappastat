package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/grsakea/kappastat/common"
	"github.com/mrshankly/go-twitch/twitch"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/redis.v3"
	"net/http"
	"strconv"
	"time"
)

func getDB() *mgo.Database {
	temp, _ := mgo.Dial("127.0.0.1")
	return (temp.DB("twitch"))
}

func getRedis() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return
}

func apiFollowing(w http.ResponseWriter, r *http.Request) {
	var ret []twitch.UserS
	db := getDB()
	db.C("follow").Find(nil).All(&ret)
	data, _ := json.Marshal(ret)
	w.Write(data)
}

func apiStat(r *http.Request, params martini.Params) string {
	var ret []kappastat.StatEntry
	var dur int
	var err error
	db := getDB()
	temp := r.URL.Query().Get("duration")
	dur, err = strconv.Atoi(temp)

	if err != nil {
		print(err)
		dur = int(1 * time.Minute.Nanoseconds())
	} else {
		dur *= int(time.Minute.Nanoseconds())
	}
	db.C("stat_entries").Find(bson.M{"channel": params["streamer"], "duration": dur}).All(&ret)
	data, _ := json.Marshal(ret)
	return string(data)
}

func apiChannel(w http.ResponseWriter, r *http.Request) {
	var ret []twitch.UserS
	db := getDB()
	db.C("follow").Find(nil).All(&ret)
	data, _ := json.Marshal(ret)
	w.Write(data)
}
