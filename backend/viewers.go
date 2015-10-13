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

	for {
		select {
		case msg := <-c:
			followed = followedHandler(followed, msg)
		case <-ticker:
			for _, v := range followed {
				infos <- fetchViewers(client, cBroadcast, v, online)
			}
		}
	}
}

func fetchViewers(client *twitch.Client, cBroadcast chan Message, chanName string, online map[string]bool) kappastat.ViewerCount {

	channel, err := client.Streams.Channel(chanName)
	if err != nil {
		channel, err = client.Streams.Channel(chanName)
		if err != nil {
			log.Print(err)
		}
	}

	if channel.Stream.Viewers != 0 && online[chanName] == false {
		online[chanName] = true
		cBroadcast <- Message{StartBroadcast, chanName}
		log.Print(chanName, " Started Broadcast")
	} else if channel.Stream.Viewers == 0 && online[chanName] == true {
		online[chanName] = false
		cBroadcast <- Message{EndBroadcast, chanName}
		log.Print(chanName, " Ended Broadcast")
	}

	return kappastat.ViewerCount{chanName, time.Now(), channel.Stream.Viewers}
}
