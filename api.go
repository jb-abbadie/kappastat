package main

import (
	"encoding/json"
	"github.com/gocraft/web"
	"github.com/grsakea/kappastat/backend"
	"gopkg.in/mgo.v2/bson"
)

func (c *Context) apiViewer(w web.ResponseWriter, r *web.Request) {
	var ret []backend.ViewerCount
	c.db.C("viewer_count").Find(bson.M{"channel": r.PathParams["streamer"]}).All(&ret)

	data, _ := json.Marshal(ret)
	w.Write(data)

}

func (c *Context) apiFollowing(w web.ResponseWriter, r *web.Request) {
	data, _ := json.Marshal(c.backend.ListStreams())
	w.Write(data)
}
