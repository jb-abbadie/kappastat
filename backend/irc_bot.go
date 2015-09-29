package main

import (
	"github.com/sorcix/irc"
	"log"
	"net"
	"time"
)

type IrcBot struct {
	server string
	login  string
	pass   string
	conn   net.Conn
	reader *irc.Decoder
	writer *irc.Encoder
	data   chan *irc.Message
}

func (b *IrcBot) loop() {
	for {
		b.conn.SetDeadline(time.Now().Add(60 * time.Second))
		msg, err := b.reader.Decode()
		if err != nil {
			log.Print("IRC channel closed :", err)
			close(b.data)
			return
		} else {
			b.data <- msg
		}
	}
}

func (b *IrcBot) reconnect() {
	var err error
	b.conn, err = net.Dial("tcp", b.server)
	backoff := 30 * time.Second
	for err != nil {
		b.conn, err = net.Dial("tcp", b.server)
		time.Sleep(backoff)
		backoff *= 2
		log.Print("Error connecting", err)
		log.Print("Retrying in ", backoff)
	}
}

func (b *IrcBot) connect() error {
	log.Print("connecting to IRC")
	b.data = make(chan *irc.Message)
	var err error
	b.conn, err = net.Dial("tcp", b.server)
	if err != nil {
		log.Print("Error connecting", err)
		return err
	}

	b.writer = irc.NewEncoder(b.conn)
	b.reader = irc.NewDecoder(b.conn)

	loginMessages := []irc.Message{
		irc.Message{
			Command: irc.PASS,
			Params:  []string{b.pass},
		},
		irc.Message{
			Command: irc.NICK,
			Params:  []string{b.login},
		},
		irc.Message{
			Command:  irc.USER,
			Params:   []string{b.login, "0", "*"},
			Trailing: b.login,
		},
	}

	for _, v := range loginMessages {
		err := b.writer.Encode(&v)
		if err != nil {
			return err
		}
	}

	go b.loop()

	return nil
}
