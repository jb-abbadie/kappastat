package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/grsakea/kappastat/common"
	"github.com/mrshankly/go-twitch/twitch"
	"gopkg.in/redis.v3"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/following.html",
	"templates/viewer.html",
	"templates/stat.html",
	"templates/index.html",
	"templates/head.inc",
	"templates/channel.html",
	"templates/header.inc"))

func launchFrontend() {
	m := martini.New()
	m.Use(martini.Static("static"))
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	r := martini.NewRouter()
	r.Get("/", indexHandler)
	r.Get("/following", followHandler)
	r.Get("/stat", statHandler)
	r.Get("/viewer", viewerHandler)
	r.Get("/add/:streamer", addHandler)
	r.Get("/del/:streamer", delHandler)
	r.Get("/api/viewer/:streamer", apiViewer)
	r.Get("/api/stat/:streamer", apiStat)
	r.Get("/api/channel/:streamer", apiStat)
	r.Get("/api/following", apiFollowing)
	db := getDB()
	m.Map(db)

	m.Action(r.Handle)
	log.Print("Started Web Server")
	m.Run()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func followHandler(w http.ResponseWriter, r *http.Request) {
	var liste []twitch.UserS
	db := getDB()
	db.C("follow").Find(nil).All(&liste)
	templates.ExecuteTemplate(w, "following.html", liste)
}

func viewerHandler(w http.ResponseWriter, r *http.Request) {
	views := []kappastat.ViewerCount{}
	templates.ExecuteTemplate(w, "viewer.html", views)
}

func statHandler(w http.ResponseWriter, r *http.Request) {
	views := []kappastat.ViewerCount{}
	templates.ExecuteTemplate(w, "stat.html", views)
}

func channelHandler(w http.ResponseWriter, r *http.Request) {
	views := []kappastat.ViewerCount{}
	templates.ExecuteTemplate(w, "channel.html", views)
}

func addHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.LPush("add", params["streamer"])
	client.Close()
	fmt.Fprintf(w, "Added %s", params["streamer"])
}

func delHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.LPush("del", params["streamer"])
	client.Close()
	fmt.Fprintf(w, "Removed %s", params["streamer"])
}
