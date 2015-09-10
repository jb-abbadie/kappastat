package backend

import (
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"time"
)

func loopViewers(client *twitch.Client, c chan Message, infos chan ViewerCount) {
	followed := []string{}
	var waitTime time.Duration

	for {
		select {
		case msg := <-c:
			if msg.s == AddStream {
				followed = append(followed, msg.v)
			} else if msg.s == Stop {
				return
			} else if msg.s == RemoveStream {
				var index int
				for i, v := range followed {
					if v == msg.v {
						index = i
					}
				}
				followed = append(followed[:index], followed[index+1:]...)
			} else {
				log.Println("Signal not handled")
			}

		default:
			if waitTime > 0 {
				time.Sleep(time.Second)
				waitTime = waitTime - time.Second
			} else {
				start := time.Now()
				for _, v := range followed {
					infos <- fetchViewers(client, v)
				}
				duration := time.Since(start)
				waitTime = time.Duration((60 - int(duration.Seconds())) * 1000000000)
				log.Println("Following : ", len(followed), " Fetch took : ", duration, " waitTime : ", waitTime)
			}
		}

	}
}

func fetchViewers(client *twitch.Client, chan_string string) ViewerCount {

	channel, err := client.Streams.Channel(chan_string)
	if err != nil {
		log.Fatal(err)
	}

	return ViewerCount{chan_string, time.Now(), channel.Stream.Viewers}
}
