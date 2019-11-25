package huobiapi

import (
	"github.com/cmdedj/huobiapi/ws"
)

/// 创建WebSocket版Market客户端
func NewAsset(accessKeyId, accessKeySecret string) (*ws.Asset, error) {
	return ws.NewAsset(accessKeyId, accessKeySecret)
}
