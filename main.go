package main

import (
	"log"
	"ytics"
)

func main() {
	log.Println("Hello")
	c := tics.SetupController()
	c.Loop()
}
