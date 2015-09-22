package main

import (
	"github.com/grsakea/kappastat/backend"
)

func main() {
	launchFrontend(launchBackend())
}

func launchBackend() *backend.Controller {
	c := backend.SetupController("twitch")
	go c.Loop()
	return c
}
