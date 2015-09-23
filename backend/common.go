package main

import "log"

func followedHandler(followed []string, msg Message) []string {
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
	} else {
		log.Println("Signal not handled")
	}
	return followed
}
