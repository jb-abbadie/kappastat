package main

import (
	"errors"
	"github.com/mrshankly/go-twitch/twitch"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"os"
	"time"
)

func (c *Controller) Loop() {
	log.Print("Start Loop")

	go loopViewers(c.twitchAPI, c.cViewer, c.infosViewer)
	go loopChat(c.cChat, c.infosChat)
	go loopStat(c.cStat, c.storage.db)

	t := time.NewTicker(time.Minute).C

	for {
		select {
		case temp, ok := <-c.infosViewer:
			if !ok {
				log.Println("InfosViewer failed")
				return
			}
			storeViewerCount(c.storage.views, temp)

		case temp, ok := <-c.infosChat:
			if !ok {
				log.Println("InfosChat failed")
				return
			}
			storeChatEntry(c.storage.chat, temp)
		case <-t:
			for c.comm.LLen("add").Val() != 0 {
				val, _ := c.comm.LPop("add").Result()
				c.AddStream(val)
			}
			for c.comm.LLen("del").Val() != 0 {
				c.AddStream(c.comm.LPop("del").String())
			}
		}
	}
	log.Println("Loop failed")
}

func SetupController(dbName string) (contr *Controller) {
	store := StorageController{
		db: setupStorage(dbName),
	}
	store.views = store.db.C("viewer_count")
	store.chat = store.db.C("chat_entries")
	store.follow = store.db.C("follow")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	contr = &Controller{
		config:      LoadConfig("config.json"),
		infosChat:   make(chan ChatEntry),
		infosViewer: make(chan ViewerCount),
		cViewer:     make(chan Message),
		cChat:       make(chan Message),
		cStat:       make(chan Message),
		tracked:     make(map[string]bool),
		comm:        client,
		storage:     store,
		twitchAPI:   twitch.NewClient(&http.Client{}),
	}

	contr.loadFollowed()

	os.Setenv("GO-TWITCH_CLIENTID", contr.config.ClientID)
	return
}

func (c *Controller) AddStream(name string) error {
	_, present := c.tracked[name]
	if present {
		log.Println("Already Following")
		return errors.New("Already Following")
	}
	log.Println("Adding", name)
	user, err := c.twitchAPI.Users.User(name)
	if err != nil {
		log.Print(err)
		return err
	}
	c.storage.follow.Insert(user)

	c.tracked[name] = true

	go func(name string) {
		c.cChat <- Message{AddStream, name}
		c.cViewer <- Message{AddStream, name}
		c.cStat <- Message{AddStream, name}
	}(name)
	log.Println("Finished adding", name)
	return nil
}

func (c *Controller) RemoveStream(name string) {
	_, present := c.tracked[name]
	if !present {
		log.Println("Not Following")
		return
	}
	log.Println("Removing ", name)

	go func(name string) {
		c.cChat <- Message{RemoveStream, name}
		c.cViewer <- Message{RemoveStream, name}
		c.cStat <- Message{RemoveStream, name}
	}(name)
	delete(c.tracked, name)
}

func (c *Controller) ListStreams() []string {
	keys := make([]string, 0, len(c.tracked))
	for k := range c.tracked {
		keys = append(keys, k)
	}
	return keys
}

func (c *Controller) loadFollowed() {
	var f []twitch.UserS
	c.storage.follow.Find(nil).All(&f)

	for _, v := range f {
		c.tracked[v.Name] = true
		go func(name string) {
			c.cChat <- Message{AddStream, name}
			c.cViewer <- Message{AddStream, name}
			c.cStat <- Message{AddStream, name}
		}(v.Name)
	}
}
