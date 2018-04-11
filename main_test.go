package main

import (
	"testing"
)

func TestAppConfig_Validate(t *testing.T) {
	config = appConfig{}
	if err := config.Validate(); err == nil {
		t.Error("empty config should not be passed")
	}
	if err := config.Validate(); err == nil {
		t.Error("dnspodtoken should not empty")
	}
	config.dnspodId = "1234"
	if err := config.Validate(); err == nil {
		t.Error("only dnspod should not be passed")
	}
	config.dnspodToken = "bfc41abad5a9852380ba15a124690ec5"
	if err := config.Validate(); err == nil {
		t.Error("only dnspodToken should not be passed")
	}
	config.subDomain = "test"
	if err := config.Validate(); err == nil {
		t.Error("only dnspodId, domain and subDomain should not be passed")
	}
	config.recordId = "123456"
	if err := config.Validate(); err == nil {
		t.Error("only domain,subDomain and recordid should not be passed")
	}
	config.domain = "example.com"
	if err := config.Validate(); err == nil {
		t.Error("internal should be required")
	}
	config.internal = -1
	if err := config.Validate(); err == nil {
		t.Error("internal should exceed than 5")
	}
	config.internal = 10
	if err := config.Validate(); err != nil {
		t.Error("all config values are valid, should pass,get error:", err)
	}
}

func TestGetPublicIP(t *testing.T) {
	ip, err := GetPublicIP()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("public ip:", ip)
}

func TestGetRecord(t *testing.T) {
	config.dnspodId = "1234"
	config.dnspodToken = "helloworld"
	config.domain = "example.com"
	config.subDomain = "test"

	recordId, ip, err := GetRecord()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if ip != "127.0.0.1" {
		t.Error("ip not equal 127.0.0.1,get:", ip)
	}
	if recordId != "123456" {
		t.Error("recordid not correct,get:", recordId)
	}
	t.Log("ip:", ip, ",recordid:", recordId)
}

func TestUpdateRecord(t *testing.T) {
	config.dnspodId = "1234"
	config.dnspodToken = "helloworld"
	config.domain = "example.com"
	config.subDomain = "test"
	config.recordId = "123456"

	if err := UpdateRecord(config.recordId, "121.40.31.121"); err != nil {
		t.Error("update record fail!error:", err.Error())
	}
}
