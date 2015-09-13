package main

import (
	"github.com/grsakea/kappastat/backend"
	"github.com/mrshankly/go-twitch/twitch"
	"net/http"
)

func main() {
	launchFrontend(launchBackend())
}

func launchBackend() *backend.Controller {
	c := backend.SetupController()
	client := twitch.NewClient(&http.Client{})
	opt := &twitch.ListOptions{
		Limit:  10,
		Offset: 0,
	}
	l, _ := client.Streams.List(opt)
	for _, item := range l.Streams {
		go c.AddStream(item.Channel.Name)
	}

	go c.Loop()
	return c
}
