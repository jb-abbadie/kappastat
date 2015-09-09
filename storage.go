package tics

import (
	"gopkg.in/mgo.v2"
	"log"
)

func setupStorage(dbName string) *mgo.Database {
	client, _ := mgo.Dial("127.0.0.1")

	return client.DB(dbName)
}

func storeChatEntry(c *mgo.Collection, ce ChatEntry) {
	err := c.Insert(ce)
	if err != nil {
		log.Println("error insert", err)
	}
	return
}

func storeViewerCount(c *mgo.Collection, vc ViewerCount) {
	err := c.Insert(vc)
	if err != nil {
		log.Println("error insert", err)
	}
	return
}
