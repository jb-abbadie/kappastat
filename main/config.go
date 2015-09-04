package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ChatLogin string
	ChatPass  string
	ClientID  string
}

func LoadConfig(path string) (ret Config) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &ret)
	if err != nil {
		log.Fatal(err)
	}

	return
}
