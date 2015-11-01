package main

import "log"

func followedHandler(followed []string, msg Message) ([]string, bool) {
	if msg.s == AddStream {
		followed = append(followed, msg.v)
	} else if msg.s == RemoveStream {
		var index int
		for i, v := range followed {
			if v == msg.v {
				index = i
			}
		}
		followed = append(followed[:index], followed[index+1:]...)
	} else if msg.s == Stop {
		return followed, false

	} else {
		log.Println("Signal not handled")
	}
	return followed, true
}
