package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/grsakea/kappastat/common"
	"github.com/mrshankly/go-twitch/twitch"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

func getDB() *mgo.Database {
	temp, _ := mgo.Dial("127.0.0.1")
	return (temp.DB("twitch"))
}

func apiViewer(w http.ResponseWriter, r *http.Request, params martini.Params) {
	var ret []kappastat.ViewerCount
	db := getDB()
	db.C("viewer_count").Find(bson.M{"channel": params["streamer"]}).All(&ret)
	data, _ := json.Marshal(ret)
	w.Write(data)

}

func apiFollowing(w http.ResponseWriter, r *http.Request) {
	var ret []twitch.UserS
	db := getDB()
	db.C("follow").Find(nil).All(&ret)
	data, _ := json.Marshal(ret)
	w.Write(data)
}

func apiStat(w http.ResponseWriter, r *http.Request, params martini.Params) {
	var ret []kappastat.StatEntry
	var dur int
	var err error
	db := getDB()
	temp := r.URL.Query().Get("duration")
	dur, err = strconv.Atoi(temp)

	if err != nil {
		dur = int(5 * time.Minute.Nanoseconds())
	}
	db.C("stat_entries").Find(bson.M{"channel": params["streamer"], "duration": dur}).All(&ret)
	data, _ := json.Marshal(ret)
	w.Write(data)
}
