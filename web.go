package main

import (
	"fmt"
	"github.com/gocraft/web"
	"github.com/grsakea/kappastat/backend"
	"gopkg.in/mgo.v2"
	"html/template"
	"log"
	"net/http"
)

type Test struct {
	Views []backend.ViewerCount
}

type Context struct {
	db      *mgo.Database
	backend *backend.Controller
}

var Backend *backend.Controller
var templates = template.Must(template.ParseFiles("templates/following.html", "templates/viewer.html"))

func launchFrontend(c *backend.Controller) {
	Backend = c
	router := web.New(Context{})
	router.Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware(web.StaticMiddleware("static")).
		Middleware((*Context).setContext)
	router.Get("/following", (*Context).followHandler)
	router.Get("/viewer", (*Context).viewerHandler)
	router.Get("/add/:streamer", (*Context).addHandler)
	router.Get("/api/viewer/:streamer", (*Context).apiViewer)
	router.Get("/api/following", (*Context).apiFollowing)

	log.Print("Started Web Server")
	log.Fatal(http.ListenAndServe("127.0.0.1:6969", router))
}

func (c *Context) setContext(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	temp, _ := mgo.Dial("127.0.0.1")
	c.db = temp.DB("twitch")
	c.backend = Backend
	next(w, r)
}

func (c *Context) followHandler(w web.ResponseWriter, r *web.Request) {
	liste := Backend.ListStreams()
	templates.ExecuteTemplate(w, "following.html", liste)
}

func (c *Context) viewerHandler(w web.ResponseWriter, r *web.Request) {
	views := []backend.ViewerCount{}

	templates.ExecuteTemplate(w, "viewer.html", views)
}

func (c *Context) addHandler(w web.ResponseWriter, r *web.Request) {
	Backend.AddStream(r.PathParams["streamer"])
	fmt.Fprintf(w, "Added %s", r.PathParams["streamer"])
}
