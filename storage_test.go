package tics

import (
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	db := setupStorage("twitch_test")

	ce := db.C("chat_entries")
	vc := db.C("viewer_count")
	ce.DropCollection()
	vc.DropCollection()

	storeChatEntry(ce, ChatEntry{"testChan", "testSender", time.Now(), "testText"})
	storeViewerCount(vc, ViewerCount{"testChan", time.Now(), 42})

	count, err := ce.Count()
	if count != 1 {
		t.Error("Chat Entry storage invalid")
	}
	if err != nil {
		t.Error("Chat Entry fetching error")
	}
	count, err = vc.Count()
	if count != 1 {
		t.Error("ViewerCount storage invalid")
	}
	if err != nil {
		t.Error("ViewerCount fetching error")
	}
}
