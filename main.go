package huobiapi

import (
	"github.com/cmdedj/huobiapi/client"
	"github.com/cmdedj/huobiapi/ws"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	// log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}

/// 创建WebSocket版资产客户端
func NewAsset(accessKeyId, accessKeySecret string) (*ws.Asset, error) {
	return ws.NewAsset(accessKeyId, accessKeySecret)
}

func NewClient(accessKeyId, accessKeySecret string) *client.Client {
	return client.NewClient(accessKeyId, accessKeySecret)
}
