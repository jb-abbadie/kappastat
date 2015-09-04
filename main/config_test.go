package main

import "testing"

func TestConfigLoading(t *testing.T) {
	conf := LoadConfig("exampleconfig.json")

	if conf.ChatLogin != "test_login" {
		t.Error("Incorect Chat Login ", conf.ChatLogin)
	}
	if conf.ChatPass != "test_pass" {
		t.Error("Incorect Chat Pass")
	}
	if conf.ClientID != "test_id" {
		t.Error("Incorect Chat ID", conf.ClientID)
	}

}
