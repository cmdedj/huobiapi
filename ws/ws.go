package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
}

var SafeWebSocketDestroyError = fmt.Errorf("connection destroy by user")

// SafeWebSocket 安全的WebSocket封装
// 保证读取和发送操作是并发安全的，支持自定义保持alive函数
type SafeWebSocket struct {
	ws              *websocket.Conn
	listener        SafeWebSocketMessageListener
	sendQueue       chan []byte
	lastError       error
	runningTaskSend bool
	runningTaskRead bool
}

type SafeWebSocketMessageListener = func(b []byte)

// NewSafeWebSocket 创建安全的WebSocket实例并连接
func NewSafeWebSocket(endpoint string) (*SafeWebSocket, error) {
	ws, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	s := &SafeWebSocket{ws: ws, sendQueue: make(chan []byte, 1000)}

	go func() {
		s.runningTaskSend = true
		for s.lastError == nil {
			b := <-s.sendQueue
			err := s.ws.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				s.lastError = err
				break
			}
		}
		s.runningTaskSend = false
	}()

	go func() {
		s.runningTaskRead = true
		for s.lastError == nil {
			_, b, err := s.ws.ReadMessage()
			if err != nil {
				s.lastError = err
				break
			}
			s.listener(b)
		}
		s.runningTaskRead = false
	}()

	return s, nil
}

// Listen 监听消息
func (s *SafeWebSocket) Listen(h SafeWebSocketMessageListener) {
	s.listener = h
}

// Send 发送消息
func (s *SafeWebSocket) Send(b []byte) {
	s.sendQueue <- b
}

// Destroy 销毁
func (s *SafeWebSocket) Destroy() (err error) {
	s.lastError = SafeWebSocketDestroyError
	for !s.runningTaskRead && !s.runningTaskSend {
		time.Sleep(time.Millisecond * 100)
	}
	if s.ws != nil {
		err = s.ws.Close()
		s.ws = nil
	}
	s.listener = nil
	s.sendQueue = nil
	return err
}

// Loop 进入事件循环，直到连接关闭才退出
func (s *SafeWebSocket) Loop() error {
	for s.lastError == nil {
		time.Sleep(time.Millisecond * 100)
	}
	return s.lastError
}
