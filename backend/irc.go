package backend

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"log"
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
			messageHandler(bot.Sender, infos, msg)
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
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{"#" + name},
	})
}

func removeChannel(s ircx.Sender, name string) {
	s.Send(&irc.Message{
		Command: irc.PART,
		Params:  []string{"#" + name},
	})
}

func messageHandler(s ircx.Sender, infos chan ChatEntry, m *irc.Message) {
	handled := make(map[string]bool)
	handled[irc.PING] = true
	handled[irc.PRIVMSG] = true
	handled[irc.RPL_WELCOME] = true
	handled[irc.RPL_YOURHOST] = true
	handled[irc.RPL_CREATED] = true
	handled[irc.RPL_MYINFO] = true
	handled[irc.RPL_MOTD] = true
	handled[irc.RPL_MOTDSTART] = true
	handled[irc.RPL_ENDOFMOTD] = true
	handled[irc.RPL_NAMREPLY] = true
	handled[irc.RPL_ENDOFNAMES] = true
	handled[irc.JOIN] = true

	if !handled[m.Command] {
		log.Println("Unhandled Message ", m.Command)
		return
	}
	if m.Command == irc.PING {
		PingHandler(s, m)
		return
	} else if m.Command == irc.PRIVMSG {
		PrivmsgHandler(s, infos, m)
	} else if m.Command == irc.JOIN {
		log.Println("Joined", m.Params[0][1:])
	}
}

func PingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}

func PrivmsgHandler(s ircx.Sender, infos chan ChatEntry, m *irc.Message) {
	channelName := m.Params[0][1:]
	sender := m.User
	text := m.Trailing
	infos <- ChatEntry{channelName, sender, time.Now(), text}
}
