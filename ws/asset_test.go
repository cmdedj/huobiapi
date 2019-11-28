package ws

import (
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func init() {
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
			_ = json
		},
	)
	log.Info("accounts sub is ", accountsSubOk)

	ordersSubOk := asset.Subscribe(
		NewOrdersSubData(),
		func(json *simplejson.Json) {
			_ = json
		},
	)
	log.Info("orders sub is ", ordersSubOk)

	time.Sleep(1 * time.Hour)
}
