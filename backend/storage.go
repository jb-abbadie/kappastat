package backend

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
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

func fetchStatData(db *mgo.Database, channel string, from time.Time, to time.Time) statData {
	vc := db.C("viewer_count").Find(bson.M{
		"channel": channel,
		"time":    bson.M{"$gt": from, "$lt": to}})
	itV := vc.Iter()
	lenV, _ := vc.Count()
	if lenV == 0 {
		vc := db.C("viewer_count").Find(bson.M{
			"channel": channel,
			"time":    bson.M{"$lt": to}}).Limit(1)
		itV = vc.Iter()
		lenV, _ = vc.Count()

		if lenV != 1 {
			log.Panic("Could not fetch stat data for ", channel)
		}
	}
	ce := db.C("chat_entries").Find(bson.M{
		"channel": channel,
		"time": bson.M{
			"$gt": from,
			"$lt": to}})
	itC := ce.Iter()
	lenC, _ := ce.Count()

	return statData{itC, lenC, itV, lenV}
}

func storeStatEntry(c *mgo.Collection, se StatEntry) {
	err := c.Insert(se)
	if err != nil {
		log.Println("error insert", err)
	}
	return
}
