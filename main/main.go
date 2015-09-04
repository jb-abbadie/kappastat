package main

import (
	"fmt"
	"github.com/mrshankly/go-twitch/twitch"
	"net/http"
	"os"
	"time"
)

func main() {

	conf := LoadConfig("config.json")

	fmt.Printf("Hello, world\n")
	os.Setenv("GO-TWITCH_CLIENTID", conf.ClientID)
	fmt.Println(os.Getenv("GO-TWITCH_CLIENTID"))
	c := make(chan Message)
	infos := make(chan StreamState)
	client := twitch.NewClient(&http.Client{})
	go loopViewers(client, c, infos)
	tracked := make(map[string]Histo)
	addStream(&tracked, c, "monstercat")
	for {
		temp := <-infos
		fmt.Println(temp)
	}
}

func addStream(tracked *map[string]Histo, c chan Message, name string) {
	_, present := (*tracked)[name]
	if present {
		fmt.Println("Already Following")
		return
	}

	(*tracked)[name] = Histo{make(map[time.Time]int), name}
	c <- Message{AddStream, name}
}
