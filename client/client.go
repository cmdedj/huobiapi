package client

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"strings"
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
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("request error")
		return nil, err
	}

	log.Info("response json ", resp.String())

	jsonData, err := simplejson.NewJson(resp.Bytes())
	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("json decode error")
		return nil, err
	}

	return jsonData, err

}

func (c *Client) GetRequest(path string, param ParamData) (*simplejson.Json, error) {
	return c.Request("GET", path, param)
}

func (c *Client) GetAccountId(accountType string) (string, error) {
	result, err := c.GetRequest("/v1/account/accounts", nil)
	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetAccountId request error")
		return "", err
	}

	var accountId string

	for _, v := range result.Get("data").MustArray() {
		data := v.(map[string]interface{})
		atType := data["type"].(string)

		if atType == accountType {
			accountId = data["id"].(json.Number).String()
		}

		return accountId, nil
	}

	return accountId, err
}

func (c *Client) GetBalance(accountId string) ([]*Balance, error) {
	result, err := c.GetRequest("/v1/account/accounts/"+accountId+"/balance", nil)

	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetBalance request error")
		return nil, err
	}

	list := result.Get("data").Get("list")
	listBytes, err := list.Encode()

	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetBalance json encode error")
		return nil, err
	}

	balances := make([]*Balance, 0)
	err = json.Unmarshal(listBytes, &balances)
	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetBalance json unmarshal error")
		return nil, err
	}

	return balances, err

}

func (c *Client) GetDepositAndWithdraw(dwType, currency, from, size, direct string) ([]*DepositAndWithdraw, error) {

	param := make(ParamData)
	param["type"] = dwType

	if currency != "" {
		param["currency"] = currency
	}
	if from != "" {
		param["from"] = from
	}
	if size != "" {
		param["size"] = size
	}
	if direct != "" {
		param["direct"] = direct
	}

	result, err := c.GetRequest("/v1/query/deposit-withdraw", param)
	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetDepositAndWithdraw request error")
		return nil, err
	}

	list := result.Get("data")
	listBytes, err := list.Encode()

	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetDepositAndWithdraw json encode error")
		return nil, err
	}

	daws := make([]*DepositAndWithdraw, 0)
	err = json.Unmarshal(listBytes, &daws)
	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetDepositAndWithdraw json unmarshal error")
		return nil, err
	}

	return daws, err
}

func (c *Client) GetOrders(symbol, states, orderTypes, startDate, endDate, from, direct, size string) ([]*Order, error) {
	param := make(ParamData)
	param["symbol"] = strings.ToLower(symbol)
	param["states"] = states

	if orderTypes != "" {
		param["types"] = orderTypes
	}
	if startDate != "" {
		param["start-date"] = startDate
	}
	if endDate != "" {
		param["end-date"] = endDate
	}
	if from != "" {
		param["from"] = from
	}
	if direct != "" {
		param["direct"] = direct
	}
	if size != "" {
		param["size"] = size
	}

	result, err := c.GetRequest("/v1/order/orders", param)
	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetOrders request error")
		return nil, err
	}

	list := result.Get("data")
	listBytes, err := list.Encode()

	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetOrders json encode error")
		return nil, err
	}

	orders := make([]*Order, 0)
	err = json.Unmarshal(listBytes, &orders)
	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetOrders json unmarshal error")
		return nil, err
	}

	return orders, err

}

func (c *Client) GetLatestSymbolPrice(symbol string) (float64, error) {
	param := make(ParamData)
	param["symbol"] = strings.ToLower(symbol)

	result, err := c.GetRequest("/market/trade", param)
	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetLatestSymbolPrice request error")
		return 0, err
	}

	price, err := result.Get("tick").Get("data").GetIndex(0).Get("price").Float64()

	if err != nil {
		log.WithFields(log.Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Error("GetLatestSymbolPrice json get error")
		return 0, err
	}

	return price, err

}
