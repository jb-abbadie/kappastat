package tics

import (
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"net/http"
	"os"
)

func (c *Controller) Loop() {

	go loopViewers(c.twitchAPI, c.cViewer, c.infosViewer)
	go loopChat(c.cChat, c.infosChat)

	for {
		select {
		case temp, ok := <-c.infosViewer:
			if !ok {
				return
			}
			storeViewerCount(c.storage.views, temp)

		case temp, ok := <-c.infosChat:
			if !ok {
				return
			}
			storeChatEntry(c.storage.chat, temp)
		default:
		}
	}
}

func SetupController() (contr *Controller) {
	store := StorageController{
		db: setupStorage("twitch"),
	}
	store.views = store.db.C("viewer_count")
	store.chat = store.db.C("chat_entries")

	contr = &Controller{
		config:      LoadConfig("config.json"),
		infosChat:   make(chan ChatEntry),
		infosViewer: make(chan ViewerCount),
		cViewer:     make(chan Message),
		cChat:       make(chan Message),
		tracked:     make(map[string]bool),
		storage:     store,
		twitchAPI:   twitch.NewClient(&http.Client{}),
	}

	os.Setenv("GO-TWITCH_CLIENTID", contr.config.ClientID)
	return
}

func (c *Controller) AddStream(name string) {
	_, present := c.tracked[name]
	if present {
		log.Println("Already Following")
		return
	}
	log.Println("Adding", name)

	c.tracked[name] = true
	c.cChat <- Message{AddStream, name}
	c.cViewer <- Message{AddStream, name}
}

func (c *Controller) RemoveStream(name string) {
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
