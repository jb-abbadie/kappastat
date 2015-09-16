package backend

import (
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"time"
)

func loopViewers(client *twitch.Client, c chan Message, infos chan ViewerCount) {
	followed := []string{}
	ticker := time.NewTicker(time.Minute).C

	for {
		select {
		case msg := <-c:
			followed = followedHandler(followed, msg)
		case <-ticker:
			for _, v := range followed {
				infos <- fetchViewers(client, v)
			}
		}
	}
}

func fetchViewers(client *twitch.Client, chan_string string) ViewerCount {

	channel, err := client.Streams.Channel(chan_string)
	if err != nil {
		channel, err = client.Streams.Channel(chan_string)
		if err != nil {
			log.Fatal(err)
		}
	}

	return ViewerCount{chan_string, time.Now(), channel.Stream.Viewers}
}
