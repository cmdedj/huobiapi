package ws

import (
	"github.com/bitly/go-simplejson"
	"github.com/globalsign/mgo/bson"
)

type pongData struct {
	Op string `json:"op"`
	Ts int64  `json:"ts"`
}

type pingData struct {
	Op string `json:"op"`
	Ts int64  `json:"ts"`
}

type subData struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}

type reqData struct {
	Req string `json:"req"`
	ID  string `json:"id"`
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
