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
	db := setupStorage("twitch")
	ce := db.C("chat_entries")
	vc := db.C("viewer_count")
	tracked := make(map[string]bool)

	infosViewer := make(chan ViewerCount)
	infosChat := make(chan ChatEntry)
	client := twitch.NewClient(&http.Client{})
	cViewer := make(chan Message)
	cChat := make(chan Message)
	go loopViewers(client, cViewer, infosViewer)
	go loopChat(cChat, infosChat)
	go addStream(&tracked, cViewer, cChat, "lirik")
	go addStream(&tracked, cViewer, cChat, "castro_1021")
	go addStream(&tracked, cViewer, cChat, "itmejp")
	go addStream(&tracked, cViewer, cChat, "jcarverpoker")
	go addStream(&tracked, cViewer, cChat, "monstercat")
	time.Sleep(1 * time.Second)
	for {
		select {
		case temp, ok := <-infosViewer:
			if !ok {
				return
			}
			storeViewerCount(vc, temp)

		case temp, ok := <-infosChat:
			if !ok {
				return
			}
			storeChatEntry(ce, temp)
		default:
		}

	}
}

func addStream(tracked *map[string]bool, cViewer chan Message, cChat chan Message, name string) {
	_, present := (*tracked)[name]
	if present {
		log.Println("Already Following")
		return
	}
	log.Println("Adding ", name)

	(*tracked)[name] = true
	cChat <- Message{AddStream, name}
	cViewer <- Message{AddStream, name}
}

func removeStream(tracked *map[string]bool, cViewer chan Message, cChat chan Message, name string) {
	_, present := (*tracked)[name]
	if !present {
		log.Println("Not Following")
		return
	}
	log.Println("Removing ", name)
	cChat <- Message{RemoveStream, name}
	cViewer <- Message{RemoveStream, name}
	delete(*tracked, name)
}
