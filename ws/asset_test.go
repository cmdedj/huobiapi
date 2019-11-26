package ws

import (
	"github.com/bitly/go-simplejson"
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


func TestAsset_Auth(t *testing.T) {

	asset, _ := NewAsset("a4382164-ed2htwf5tf-6d55e15e-701e5", "e7de9097-0adeb442-66b6f2d7-76752")
	ok := asset.Auth()
	log.Info("auth is ", ok)

	accountsSubOk := asset.Subscribe(
		NewAccountsSubData("0"),
		func(json *simplejson.Json) {
			log.Debug(json)
		},
	)
	log.Info("accounts sub is ", accountsSubOk)

	ordersSubOk := asset.Subscribe(
		NewOrdersSubData(),
		func(json *simplejson.Json) {
			log.Debug(json)
		},
	)
	log.Info("orders sub is ", ordersSubOk)

	asset.Loop()

}
