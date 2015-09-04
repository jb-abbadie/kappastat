package main

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"log"
	"strings"
	"time"
)

func setupChat() *ircx.Bot {
	conf := LoadConfig("config.json")
	bot := ircx.WithLogin("irc.twitch.tv:6667", conf.ChatLogin, conf.ChatLogin, conf.ChatPass)
	err := bot.Connect()
	if err != nil {
		log.Panicln(err)
	}
	return bot
}

func loopChat(c chan Message, infos chan ChatEntry) {
	bot := setupChat()
	for {
		select {
		case msg, ok := <-bot.Data:
			if !ok {
				return
			}
			messageHandler(bot.Sender, msg)
		case msg, ok := <-c:
			if !ok {
				return
			}
			if msg.s == AddStream {
				addChannel(bot.Sender, msg.v)
			}
		}
	}
}

func addChannel(s ircx.Sender, name string) {
	log.Println("#" + name)

	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{"#" + name},
	})
}

func messageHandler(s ircx.Sender, m *irc.Message) {
	if m.Command != irc.PRIVMSG {
		log.Println("Unhandled Message ", m.Command)
		return
	}

	split := strings.Split(m.String(), " ")

	if strings.Contains(split[0], "twitchnotify") {
		log.Println("Subscriber")
		log.Println(m)
	} else {
		channelName := split[2][1:]
		sender := split[0]
		sender = sender[1:strings.IndexRune(sender, '!')]
		text := strings.Join(split[3:], " ")[1:]
		entry := ChatEntry{channelName, sender, time.Now(), text}
		log.Println(entry)

	}

}
