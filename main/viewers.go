package main

import (
	"fmt"
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"time"
)

func loopViewers(client *twitch.Client, c chan Message, infos chan StreamState) {
	followed := []string{}

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
			start := time.Now()
			for _, v := range followed {
				infos <- fetchViewers(client, v)
			}
			duration := time.Since(start)
			waitTime := time.Duration(60 - int(duration.Seconds()))
			print(waitTime)
			fmt.Println("Following : ", len(followed), " Fetch took : ", duration)
			time.Sleep(waitTime * time.Second)
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
