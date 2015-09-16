package main

import (
	"fmt"
	"github.com/gocraft/web"
	"github.com/grsakea/kappastat/backend"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
)

type Test struct {
	Views []backend.ViewerCount
}

type Context struct {
	db *mgo.Database
}

var Backend *backend.Controller
var templates = template.Must(template.ParseFiles("templates/following.html", "templates/viewer.html"))

func launchFrontend(c *backend.Controller) {
	Backend = c
	router := web.New(Context{})
	router.Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*Context).setContext)
	router.Get("/following", (*Context).followHandler)
	router.Get("/viewer/:streamer", (*Context).viewerHandler)
	router.Get("/add/:streamer", (*Context).addHandler)

	log.Print("Started Web Server")
	log.Fatal(http.ListenAndServe("127.0.0.1:6969", router))
}

func (c *Context) setContext(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	temp, _ := mgo.Dial("127.0.0.1")
	c.db = temp.DB("twitch")
	next(w, r)
}

func (c *Context) followHandler(w web.ResponseWriter, r *web.Request) {
	liste := Backend.ListStreams()
	templates.ExecuteTemplate(w, "templates/following.html", liste)
}

func (c *Context) viewerHandler(w web.ResponseWriter, r *web.Request) {
	coll := c.db.C("viewer_count")
	views := []backend.ViewerCount{}
	streamer := r.PathParams["streamer"]
	coll.Find(bson.M{"channel": streamer}).All(&views)

	templates.ExecuteTemplate(w, "viewer.html", views)
}

func (c *Context) addHandler(w web.ResponseWriter, r *web.Request) {
	Backend.AddStream(r.PathParams["streamer"])
	fmt.Fprintf(w, "Added %s", r.PathParams["streamer"])
}
