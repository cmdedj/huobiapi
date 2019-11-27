package client

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

var client *Client

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	// log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	client = NewClient("a4382164-ed2htwf5tf-6d55e15e-701e5", "e7de9097-0adeb442-66b6f2d7-76752")
}

func TestClient_GetAccountId(t *testing.T) {
	aid, err := client.GetAccountId(AccountTypeSpot)
	log.Println("aid", aid)
	if err != nil {
		log.Error(err)
	}
}

func TestClient_GetBalance(t *testing.T) {
	aid, _ := client.GetAccountId(AccountTypeSpot)
	balance, err := client.GetBalance(aid)
	log.Println("balance", balance)
	if err != nil {
		log.Error(err)
	}
}
