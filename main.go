package main

import (
	"encoding/json"
	"github.com/grsakea/kappastat/backend"
	"github.com/mrshankly/go-twitch/twitch"
	"io/ioutil"
	"net/http"
)

func main() {
	launchFrontend(launchBackend())
}

func launchBackend() *backend.Controller {
	c := backend.SetupController()
	client := twitch.NewClient(&http.Client{})

	data, err := ioutil.ReadFile("following")
	if err == nil {
		loadFollowing(c, data)
		print(data)
	} else {
		opt := &twitch.ListOptions{
			Limit:  10,
			Offset: 0,
		}
		l, _ := client.Streams.List(opt)
		for _, item := range l.Streams {
			go c.AddStream(item.Channel.Name)
		}
	}
	go c.Loop()
	return c
}

func loadFollowing(c *backend.Controller, data []byte) {
	var f []string
	json.Unmarshal(data, &f)

	for i := range f {
		go c.AddStream(f[i])
	}
}
