package client

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	// log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}


func TestClient_GetAccountId(t *testing.T) {
	client := NewClient("a4382164-ed2htwf5tf-6d55e15e-701e5", "e7de9097-0adeb442-66b6f2d7-76752")
	client.GetAccountId()
}