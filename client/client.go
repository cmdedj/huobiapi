package client

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	AccessKeyId     string
	AccessKeySecret string
}

/// 全局API
const Endpoint = "https://api-aws.huobi.pro"

/// 创建新客户端
func NewClient(accessKeyId, accessKeySecret string) *Client {

	return &Client{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}

}

/// 发送请求
func (c *Client) Request(method string, path string, param ParamData) (*simplejson.Json, error) {
	if param == nil {
		param = make(ParamData)
	}

	param["AccessKeyId"] = c.AccessKeyId
	param["SignatureMethod"] = "HmacSHA256"
	param["SignatureVersion"] = "2"
	param["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	param["Signature"] = GenSignature(method, path, param, c.AccessKeySecret)

	ro := &grequests.RequestOptions{
		Params: param,
	}

	resp, err := grequests.Get(Endpoint+path, ro)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("response json ", resp.String())

	jsonData, err := simplejson.NewJson(resp.Bytes())
	if err != nil {
		log.Error(err)
	}

	return jsonData, err

}

func (c *Client) GetRequest(path string, param ParamData) (*simplejson.Json, error) {
	return c.Request("GET", path, param)
}

func (c *Client) GetAccountId(accountType string) (string, error) {
	result, err := c.GetRequest("/v1/account/accounts", nil)
	if err != nil {
		log.Error(err)
	}

	var accountId string

	for _, v := range result.Get("data").MustArray() {
		data := v.(map[string]interface{})
		accountType := data["type"].(string)

		if accountType == accountType {
			accountId = data["id"].(json.Number).String()
		}

	}

	return accountId, err
}

func (c *Client) GetBalance(accountId string) ([]*Balance, error) {
	result, err := c.GetRequest("/v1/account/accounts/"+accountId+"/balance", nil)

	if err != nil {
		log.Error(err)
	}

	list := result.Get("data").Get("list")
	listBytes, err := list.Encode()

	if err != nil {
		log.Error(err)
	}

	balances := make([]*Balance, 0)
	err = json.Unmarshal(listBytes, &balances)
	if err != nil {
		log.Error(err)
	}

	return balances, err

}
