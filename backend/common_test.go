package main

import "testing"

func TestFollowedHander(t *testing.T) {
	followed := []string{}

	followed, _ = followedHandler(followed, Message{AddStream, "test"})
	if followed[0] != "test" {
		t.Error("Error adding stream")
	}
	followed, _ = followedHandler(followed, Message{RemoveStream, "test"})
	if len(followed) != 0 {
		t.Error("Error removing stream")
	}
	_, test := followedHandler(followed, Message{Stop, ""})

	if test != false {
		t.Error("Error Stopping loop")
	}
}
