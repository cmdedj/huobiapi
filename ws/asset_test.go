package ws

import (
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestAsset_Auth(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	// log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	asset, _ := NewAsset("a4382164-ed2htwf5tf-6d55e15e-701e5", "e7de9097-0adeb442-66b6f2d7-76752")
	ok := asset.Auth()
	log.Info("auth is ", ok)

	subok := asset.Subscribe(
		NewAccountsSubData("0"),
		func(json *simplejson.Json) {
			log.Debug(json)
		},
	)
	log.Info("sub is ", subok)

	asset.Loop()

}
