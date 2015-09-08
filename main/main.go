package main

import (
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"net/http"
	"os"
)

func main() {

	client := twitch.NewClient(&http.Client{})
	db := setupStorage("twitch")
	ce := db.C("chat_entries")
	vc := db.C("viewer_count")
	c := setupController()

	go loopViewers(client, c.cViewer, c.infosViewer)
	go loopChat(c.cChat, c.infosChat)

	for {
		select {
		case temp, ok := <-c.infosViewer:
			if !ok {
				return
			}
			storeViewerCount(vc, temp)

		case temp, ok := <-c.infosChat:
			if !ok {
				return
			}
			storeChatEntry(ce, temp)
		default:
		}
	}
}

func setupController() (contr *Controller) {

	contr = &Controller{
		config:      LoadConfig("config.json"),
		infosChat:   make(chan ChatEntry),
		infosViewer: make(chan ViewerCount),
		cViewer:     make(chan Message),
		cChat:       make(chan Message),
		tracked:     make(map[string]bool),
	}

	os.Setenv("GO-TWITCH_CLIENTID", contr.config.ClientID)
	return
}

func (c *Controller) addStream(name string) {
	_, present := c.tracked[name]
	if present {
		log.Println("Already Following")
		return
	}
	log.Println("Adding ", name)

	c.tracked[name] = true
	c.cChat <- Message{AddStream, name}
	c.cViewer <- Message{AddStream, name}
}

func (c *Controller) removeStream(name string) {
	_, present := c.tracked[name]
	if !present {
		log.Println("Not Following")
		return
	}
	log.Println("Removing ", name)
	c.cChat <- Message{RemoveStream, name}
	c.cViewer <- Message{RemoveStream, name}
	delete(c.tracked, name)
}
