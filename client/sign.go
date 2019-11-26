package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

func getMapKeys(m map[string]string) (keys []string) {
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func sortKeys(keys []string) []string {
	sort.Strings(keys)
	return keys
}

func computeHmac256(data string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

/// 拼接query字符串
func encodeQueryString(query map[string]string) string {
	var keys = sortKeys(getMapKeys(query))
	var keysLen = len(keys)
	var lines = make([]string, keysLen)
	for i := 0; i < keysLen; i++ {
		var k = keys[i]
		lines[i] = url.QueryEscape(k) + "=" + url.QueryEscape(query[k])
	}
	return strings.Join(lines, "&")
}

func GenSignature(method string, path string, data ParamData, accessKeySecret string) string {
	method = strings.ToUpper(method)
	var pre = method + "\n" + "api-aws.huobi.pro" + "\n" + path + "\n"
	eqs := encodeQueryString(data)
	eqs = pre + eqs
	return computeHmac256(eqs, accessKeySecret)
}
