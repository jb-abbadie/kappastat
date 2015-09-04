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
	infosViewer := make(chan StreamState)
	infosChat := make(chan ChatEntry)
	client := twitch.NewClient(&http.Client{})
	cViewer := make(chan Message)
	cChat := make(chan Message)
	go loopViewers(client, cViewer, infosViewer)
	go loopChat(cChat, infosChat)
	tracked := make(map[string]Histo)
	addStream(&tracked, cViewer, cChat, "monstercat")
	addStream(&tracked, cViewer, cChat, "lirik")
	for {
		temp := <-infosViewer
		fmt.Println(temp)
	}
}

func addStream(tracked *map[string]Histo, cViewer chan Message, cChat chan Message, name string) {
	_, present := (*tracked)[name]
	if present {
		fmt.Println("Already Following")
		return
	}

	(*tracked)[name] = Histo{make(map[time.Time]int), name}
	cChat <- Message{AddStream, name}
	cViewer <- Message{AddStream, name}
}
