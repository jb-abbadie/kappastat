package main

import (
	"github.com/grsakea/kappastat/backend"
)

func main() {
	launchBackend()
	launchFrontend()
}

func launchBackend() *backend.Controller {
	c := backend.SetupController("twitch")
	go c.Loop()
	return c
}
