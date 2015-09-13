package backend

import "testing"

func TestFollowedHander(t *testing.T) {
	followed := []string{}

	followed = followedHandler(followed, Message{AddStream, "test"})
	if followed[0] != "test" {
		t.Error("Error adding stream")
	}
	followed = followedHandler(followed, Message{RemoveStream, "test"})
	if len(followed) != 0 {
		t.Error("Error removing stream")
	}
}
