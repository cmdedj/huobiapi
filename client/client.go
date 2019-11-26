package client

import (
	"github.com/bitly/go-simplejson"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	AccessKeyId      string
	AccessKeySecret  string
}

/// 全局API
const Endpoint = "https://api-aws.huobi.pro"

/// 创建新客户端
func NewClient(accessKeyId, accessKeySecret string) *Client {

	client := &Client{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	return client
}

type ParamData = map[string]string

/// 发送请求
func (c *Client) Request(method string, path string, data ParamData) (*simplejson.Json, error) {
	if data == nil {
		data = make(ParamData)
	}

	data["AccessKeyId"] = c.AccessKeyId
	data["SignatureMethod"] = "HmacSHA256"
	data["SignatureVersion"] = "2"
	data["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	data["Signature"] = GenSignature(method, path, data, c.AccessKeySecret)

	ro := &grequests.RequestOptions{
		Params: data,
	}

	resp, err := grequests.Get(Endpoint, ro)
	if err != nil {
		log.Error(err)
	}

	log.Info(resp.String())

	return nil, nil

}

func (c *Client) GetRequest(path string, data ParamData) (*simplejson.Json, error) {
	return c.Request("GET", path, data)
}

func (c *Client) GetAccountId() (*simplejson.Json, error) {
	return c.GetRequest("/v1/account/accounts", nil)
}
