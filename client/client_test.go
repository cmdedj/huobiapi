package client

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

var client *Client

func init() {
	log.SetReportCaller(false)
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

func TestClient_GetDepositAndWithdraw(t *testing.T) {
	daws, err := client.GetDepositAndWithdraw(Deposit, "", "", "500", DirectDesc)

	log.WithFields(log.Fields{
		"daws": fmt.Sprintf("%+v", daws),
	}).Info("daws")

	if err != nil {
		log.Error(err)
	}

}

func TestClient_GetOrders(t *testing.T) {
	orders, err := client.GetOrders("BTCUSDT", "filled,partial-filled", "", "", "", "", DirectDesc, "100")

	log.WithFields(log.Fields{
		"orders": fmt.Sprintf("%+v", *orders[0]),
	}).Info("orders")

	if err != nil {
		log.Error(err)
	}
}

func TestClient_GetLatestSymbolPrice(t *testing.T) {
	price, err := client.GetLatestSymbolPrice("btcusdt")

	log.WithFields(log.Fields{
		"price": fmt.Sprintf("%+v", price),
	}).Info("price")

	if err != nil {
		log.Error(err)
	}
}
