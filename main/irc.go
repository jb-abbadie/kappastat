package main

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"log"
	"strings"
)

func setupChat() *ircx.Bot {
	conf := LoadConfig("config.json")
	bot := ircx.WithLogin("irc.twitch.tv:6667", conf.ChatLogin, conf.ChatLogin, conf.ChatPass)
	err := bot.Connect()
	if err != nil {
		log.Panicln(err)
	}
	bot.HandleFunc(irc.RPL_WELCOME, RegisterConnect)
	bot.HandleFunc(irc.PRIVMSG, MessageHandler)
	return bot
}

func loopChat(c chan Message) {
	bot := setupChat()
	log.Println("Start")
	bot.HandleLoop()
	log.Println("Finish")
}

func RegisterConnect(s ircx.Sender, m *irc.Message) {
	log.Println("Joined")
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{"#lirik"},
	})
}

func MessageHandler(s ircx.Sender, m *irc.Message) {
	split := strings.Split(m.String(), " ")
	if strings.Contains(split[0], "twitchnotify") {
		log.Println("Subscriber")
		log.Println(m)
	}
}
