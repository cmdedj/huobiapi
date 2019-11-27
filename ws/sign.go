package ws

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
)

// unGzipData 解压gzip的数据
func unGzipData(buf []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

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

func GenSignature(query map[string]string, accessKeySecret string) string {
	var pre = "GET" + "\n" + "api-aws.huobi.pro" + "\n" + "/ws/v1" + "\n"
	eqs := encodeQueryString(query)
	eqs = pre + eqs
	return computeHmac256(eqs, accessKeySecret)
}
