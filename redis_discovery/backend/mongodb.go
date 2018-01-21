package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
)

type Entry struct {
	Timestamp time.Time
	Text      string
}

func getMongoDB() (*mgo.Database, error) {
	service := env.GetService(mongoDbServiceInstance)
	session, err := mgo.Dial(fmt.Sprintf("%v", service.Credentials["url"]))
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB(fmt.Sprintf("%v", service.Credentials["db"])), nil
}

func insertEntry(text string) error {
	db, err := getMongoDB()
	if err != nil {
		return err
	}
	defer db.Session.Close()

	entry := Entry{
		Timestamp: time.Now(),
		Text:      text,
	}
	e := db.C("entries")
	if err := e.Insert(entry); err != nil {
		return err
	}
	return nil
}

func getEntries() ([]Entry, error) {
	db, err := getMongoDB()
	if err != nil {
		return nil, err
	}
	defer db.Session.Close()

	var result []Entry
	e := db.C("entries")
	if err := e.Find(nil).Sort("-timestamp").All(&result); err != nil {
		return nil, err
	}
	return result, nil
}
