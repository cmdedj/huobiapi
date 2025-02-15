package ws

import (
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)
}

func TestAsset_Auth(t *testing.T) {

	asset, err := NewAsset("a4382164-ed2htwf5tf-6d55e15e-701e5", "e7de9097-0adeb442-66b6f2d7-76752")
	if err != nil {
		log.Error(err)
	}

	ok := asset.Auth()
	log.Info("auth is ", ok)

	accountsSubOk := asset.Subscribe(
		NewAccountsSubData(BalanceAll),
		func(json *simplejson.Json) {
			log.Info("get balance all")
			_ = json
		},
	)
	log.Info("accounts sub is ", accountsSubOk)

	ordersSubOk := asset.Subscribe(
		NewOrdersSubData(),
		func(json *simplejson.Json) {
			log.Info("get orders")
			_ = json
		},
	)
	log.Info("orders sub is ", ordersSubOk)

	//time.Sleep(time.Second)

	//ordersUnSubOk := asset.UnSubscribe("orders.*.update")
	//log.Info("orders unsub is ", ordersUnSubOk)

	time.Sleep(1 * time.Hour)
}
