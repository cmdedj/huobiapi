package ws

import (
	"github.com/bitly/go-simplejson"
	"github.com/globalsign/mgo/bson"
)

type pingData struct {
	Op string `json:"op"`
	Ts int64  `json:"ts"`
}

type pongData struct {
	Op string `json:"op"`
	Ts int64  `json:"ts"`
}

type AuthData struct {
	Op               string `json:"op"`
	Cid              string `json:"cid"`
	AccessKeyId      string `json:"AccessKeyId"`
	SignatureMethod  string `json:"SignatureMethod"`
	SignatureVersion string `json:"SignatureVersion"`
	Timestamp        string `json:"Timestamp"`
	Signature        string `json:"Signature"`
}

type AccountsList struct {
	Op    string `json:"op"`
	Cid   string `json:"cid"`
	Topic string `json:"topic"`
}

type jsonChan = chan *simplejson.Json

// Listener 订阅事件监听器
type Listener = func(json *simplejson.Json)

type SubData interface {
	GetCid() string
	GetTopic() string
}

type AccountsSubData struct {
	Op    string `json:"op"`
	Cid   string `json:"cid"`
	Topic string `json:"topic"`
	Model string `json:"model"`
}

func (asd *AccountsSubData) GetCid() string {
	return asd.Cid
}

func (asd *AccountsSubData) GetTopic() string {
	return asd.Topic
}

func NewAccountsSubData(model string) SubData {
	return &AccountsSubData{
		Op:    "sub",
		Cid:   bson.NewObjectId().Hex(),
		Topic: "accounts",
		Model: model,
	}
}

type OrdersSubData struct {
	Op    string `json:"op"`
	Cid   string `json:"cid"`
	Topic string `json:"topic"`
}

func (osd *OrdersSubData) GetCid() string {
	return osd.Cid
}

func (osd *OrdersSubData) GetTopic() string {
	return osd.Topic
}

func NewOrdersSubData() SubData {
	return &OrdersSubData{
		Op:    "sub",
		Cid:   bson.NewObjectId().Hex(),
		Topic: "orders.*.update",
	}
}
