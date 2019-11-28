package ws

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Endpoint 行情的Websocket入口
var assetEndpoint = "wss://api-aws.huobi.pro/ws/v1"

type Asset struct {
	ws *SafeWebSocket

	listeners       map[string]Listener
	listenerMutex   sync.Mutex
	subscribedTopic map[string]SubData
	requestResultCb map[string]jsonChan

	// 掉线后是否自动重连，如果用户主动执行Close()则不自动重连
	autoReconnect bool
	autoReconnectCount int

	AccessKeyId     string
	AccessKeySecret string
}

// NewMarket 创建Market实例
func NewAsset(accessKeyId, accessKeySecret string) (asset *Asset, err error) {
	asset = &Asset{
		ws:              nil,
		autoReconnect:   true,
		autoReconnectCount: 3,
		listeners:       make(map[string]Listener),
		requestResultCb: make(map[string]jsonChan),
		subscribedTopic: make(map[string]SubData),
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}

	if err := asset.connect(); err != nil {
		return nil, err
	}

	go asset.Loop()

	return asset, err
}

// connect 连接
func (asset *Asset) connect() error {
	ws, err := NewSafeWebSocket(assetEndpoint)
	if err != nil {
		return err
	}
	asset.ws = ws
	asset.handleMessageLoop()

	return err
}

// reconnect 重新连接
func (asset *Asset) reconnect() error {

	time.Sleep(time.Second)

	err := asset.connect()
	if err != nil {

		return err
	}

	// 重新订阅
	asset.listenerMutex.Lock()
	var listeners = make(map[string]Listener)
	for k, v := range asset.listeners {
		listeners[k] = v
	}
	asset.listenerMutex.Unlock()

	for topic, listener := range listeners {
		delete(asset.subscribedTopic, topic)
		asset.Subscribe(asset.subscribedTopic[topic], listener)
	}
	return err
}

// sendMessage 发送消息
func (asset *Asset) SendMessage(message interface{}) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"message": fmt.Sprintf("%+v", message),
	}).Info("send message")
	asset.ws.Send(b)
	return nil
}

// handleMessageLoop 处理消息循环
func (asset *Asset) handleMessageLoop() {
	asset.ws.Listen(func(buf []byte) {
		msg, err := unGzipData(buf)
		if err != nil {
			log.WithFields(log.Fields{
				"err": fmt.Sprintf("%+v", err),
			}).Error("un gzip error")
			return
		}

		log.Info("response json ", string(msg))

		jsonData, err := simplejson.NewJson(msg)
		if err != nil {
			log.WithFields(log.Fields{
				"err": fmt.Sprintf("%+v", err),
			}).Error("json decode")
			return
		}

		op := jsonData.Get("op").MustString()

		// 处理ping
		if op == "ping" {
			ts := jsonData.Get("ts").MustInt64()
			err := asset.handlePing(pingData{
				Op: "ping",
				Ts: ts,
			})
			if err != nil {
				log.WithFields(log.Fields{
					"err": fmt.Sprintf("%+v", err),
				}).Error("handle ping")
			}

			return
		} else if op == "notify" {

			topic := jsonData.Get("topic").MustString()
			asset.listenerMutex.Lock()
			listener, ok := asset.listeners[topic]
			asset.listenerMutex.Unlock()
			if ok {
				listener(jsonData)
			}
			return
		} else {
			cid, _ := jsonData.Get("cid").String()
			c, ok := asset.requestResultCb[cid]
			if ok {
				c <- jsonData
			}
			return
		}
	})
}

// handlePing 处理Ping
func (asset *Asset) handlePing(ping pingData) (err error) {

	var pong = pongData{
		Op: "pong",
		Ts: ping.Ts,
	}
	err = asset.SendMessage(pong)
	if err != nil {
		return err
	}
	return nil
}

// Subscribe 订阅
func (asset *Asset) Subscribe(subData SubData, listener Listener) bool {

	// 如果未曾发送过订阅指令，则发送，并等待订阅操作结果，否则直接返回
	if _, ok := asset.subscribedTopic[subData.GetTopic()]; !ok {

		jsonData, err := asset.Request(subData.GetCid(), subData)
		if err != nil{
			log.WithFields(log.Fields{
				"err": fmt.Sprintf("%+v", err),
			}).Info("Subscribe send message error")
			return false
		}

		if jsonData != nil {
			errCode := jsonData.Get("err-code").MustInt()
			if errCode == 0 {
				return true
			} else {
				return false
			}
		}

	}

	asset.listenerMutex.Lock()
	asset.listeners[subData.GetTopic()] = listener
	asset.listenerMutex.Unlock()
	asset.subscribedTopic[subData.GetTopic()] = subData

	return true

}

// Unsubscribe 取消订阅
func (asset *Asset) UnSubscribe(topic string) bool {

	cid := bson.NewObjectId().Hex()
	jsonData, err := asset.Request(cid, UnSubData{
		Unsub: topic,
		Id:    cid,
	})
	if err != nil{
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Info("UnSubscribe send message error")
		return false
	}

	if jsonData != nil {
		errCode := jsonData.Get("err-code").MustInt()
		if errCode == 0 {
			return true
		} else {
			return false
		}
	}

	asset.listenerMutex.Lock()
	delete(asset.listeners, topic)
	asset.listenerMutex.Unlock()

	return true
}

func (asset *Asset) Request(cid string, data interface{}) (*simplejson.Json, error) {

	asset.requestResultCb[cid] = make(jsonChan)

	err := asset.SendMessage(data)
	if err != nil {
		return nil, err
	}

	var jsonData = <-asset.requestResultCb[cid]

	delete(asset.requestResultCb, cid)

	return jsonData, err
}

// Loop 进入循环
func (asset *Asset) Loop() {
	for {
		err := asset.ws.Loop()
		if err != nil {

			if err == SafeWebSocketDestroyError {
				break
			} else if asset.autoReconnect {
				err = asset.reconnect()
				if err != nil{
					asset.autoReconnectCount -= 1
				}
				if asset.autoReconnectCount < 0 {
					break
				}
				time.Sleep(3 * time.Second)
			} else {
				break
			}
		}
	}

}

// ReConnect 重新连接
func (asset *Asset) ReConnect() (err error) {

	asset.autoReconnect = true
	if err = asset.ws.Destroy(); err != nil {
		return err
	}
	return asset.reconnect()
}

// Close 关闭连接
func (asset *Asset) Close() error {

	asset.autoReconnect = false
	if err := asset.ws.Destroy(); err != nil {
		return err
	}
	return nil
}

func (asset *Asset) Auth() bool {
	params := make(map[string]string)

	params["AccessKeyId"] = asset.AccessKeyId
	params["SignatureMethod"] = "HmacSHA256"
	params["SignatureVersion"] = "2"
	params["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	cid := bson.NewObjectId().Hex()
	authData := AuthData{
		Op:               "auth",
		Cid:              cid,
		AccessKeyId:      params["AccessKeyId"],
		SignatureMethod:  params["SignatureMethod"],
		SignatureVersion: params["SignatureVersion"],
		Timestamp:        params["Timestamp"],
		Signature:        GenSignature(params, asset.AccessKeySecret),
	}

	jsonData, err := asset.Request(cid, authData)

	if err != nil{
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Info("Auth request error")
		return false
	}

	if jsonData != nil {
		errCode := jsonData.Get("err-code").MustInt()
		if errCode == 0 {
			return true
		} else {
			return false
		}
	}

	return false
}
