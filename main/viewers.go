package main

import (
	"fmt"
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"time"
)

func loopViewers(client *twitch.Client, c chan Message, infos chan StreamState) {
	followed := []string{}
	var waitTime time.Duration

	for {
		select {
		case msg := <-c:
			if msg.s == AddStream {
				followed = append(followed, msg.v)
			} else if msg.s == Stop {
				return
			} else {
				fmt.Println("Signal not handled")
			}

		default:
			if waitTime > 0 {
				time.Sleep(time.Second)
			} else {
				start := time.Now()
				for _, v := range followed {
					infos <- fetchViewers(client, v)
				}
				duration := time.Since(start)
				waitTime = time.Duration(60 - int(duration.Seconds()))
				fmt.Println("Following : ", len(followed), " Fetch took : ", duration)
			}
		}

	}
}

func fetchViewers(client *twitch.Client, chan_string string) StreamState {

	channel, err := client.Streams.Channel(chan_string)
	if err != nil {
		log.Fatal(err)
	}

	return StreamState{chan_string, time.Now(), channel.Stream.Viewers}
}
