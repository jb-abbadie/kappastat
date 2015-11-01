package main

import (
	"github.com/grsakea/kappastat/common"
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"time"
)

func loopViewers(client *twitch.Client, c chan Message, cBroadcast chan Message, infos chan kappastat.ViewerCount) {
	followed := []string{}
	online := make(map[string]bool)
	ticker := time.NewTicker(time.Minute).C
	loop := true

	for loop {
		select {
		case msg := <-c:
			followed, loop = followedHandler(followed, msg)
		case <-ticker:
			for _, v := range followed {
				infos <- fetchViewers(client, cBroadcast, v, online)
			}
		}
	}
	log.Print("Viewer loop stopped")
}

func fetchViewers(client *twitch.Client, cBroadcast chan Message, chanName string, online map[string]bool) kappastat.ViewerCount {

	channel, err := client.Streams.Channel(chanName)
	if err != nil {
		channel, err = client.Streams.Channel(chanName)
		if err != nil {
			log.Print("Error fetching viewer : ", err)
		}
	}

	if channel.Stream.Viewers != 0 && online[chanName] == false {
		online[chanName] = true
		cBroadcast <- Message{StartBroadcast, chanName}
	} else if channel.Stream.Viewers == 0 && online[chanName] == true {
		online[chanName] = false
		cBroadcast <- Message{EndBroadcast, chanName}
	}

	return kappastat.ViewerCount{chanName, time.Now(), channel.Stream.Viewers}
}
