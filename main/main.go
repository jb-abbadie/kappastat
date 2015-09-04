package main

import (
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	conf := LoadConfig("config.json")

	os.Setenv("GO-TWITCH_CLIENTID", conf.ClientID)
	infosViewer := make(chan StreamState)
	infosChat := make(chan ChatEntry)
	client := twitch.NewClient(&http.Client{})
	cViewer := make(chan Message)
	cChat := make(chan Message)
	go loopViewers(client, cViewer, infosViewer)
	go loopChat(cChat, infosChat)
	time.Sleep(2 * time.Second)
	tracked := make(map[string]Histo)
	addStream(&tracked, cViewer, cChat, "lirik")
	for {
		select {
		case temp, ok := <-infosViewer:
			if !ok {
				return
			}
			log.Println(temp)

		case temp, ok := <-infosChat:
			if !ok {
				return
			}
			log.Println(temp)
		default:
		}

	}
}

func addStream(tracked *map[string]Histo, cViewer chan Message, cChat chan Message, name string) {
	_, present := (*tracked)[name]
	if present {
		log.Println("Already Following")
		return
	}
	log.Println("Adding ", name)

	(*tracked)[name] = Histo{make(map[time.Time]int), name}
	cChat <- Message{AddStream, name}
	cViewer <- Message{AddStream, name}
}
