package bybit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/iancoleman/orderedmap"
)

type authParam struct {
	apiKey    string
	apiSecret string
}

func (a *authParam) sign(r RequestType) (string, int64, error) {
	q, ts, err := a.composeQuery(r)
	if err != nil {
		return "", 0, err
	}
	h := hmac.New(sha256.New, []byte(a.apiSecret))
	if _, err := io.WriteString(h, q); err != nil {
		return "", 0, err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), ts, nil
}

func (a *authParam) composeQuery(r RequestType) (string, int64, error) {
	sb, err := json.Marshal(r)
	if err != nil {
		return "", 0, err
	}
	var srcMap map[string]interface{}
	if err := json.Unmarshal(sb, &srcMap); err != nil {
		return "", 0, err
	}

	srcMap["api_key"] = a.apiKey
	timestamp := time.Now().UnixMilli()
	srcMap["timestamp"] = timestamp

	var keys []string
	for k := range srcMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// url.Values shouldn't work because it's map, doesn't care the order.
	dest := ""
	for _, k := range keys {
		val := ""
		switch reflect.TypeOf(srcMap[k]).Kind() {
		case reflect.Float64:
			val = strconv.FormatFloat(srcMap[k].(float64), 'f', -1, 64)
		case reflect.Int64:
			val = strconv.FormatInt(srcMap[k].(int64), 10)
		case reflect.Int:
			val = strconv.Itoa(srcMap[k].(int))
		case reflect.Bool:
			val = strconv.FormatBool(srcMap[k].(bool))
		case reflect.String:
			val = srcMap[k].(string)
		}
		dest += fmt.Sprintf("%s=%s&", k, val)
	}

	return dest[0 : len(dest)-1], timestamp, nil
}

func (a *authParam) composeBody(req RequestType) ([]byte, error) {
	sb, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var srcMap map[string]interface{}
	if err := json.Unmarshal(sb, &srcMap); err != nil {
		return nil, err
	}

	sign, timestamp, err := a.sign(req)
	if err != nil {
		return nil, err
	}

	o := orderedmap.New()
	for k, v := range srcMap {
		o.Set(k, v)
	}
	o.Set("api_key", a.apiKey)
	o.Set("timestamp", timestamp)
	o.SortKeys(sort.Strings)
	o.Set("sign", sign)

	return o.MarshalJSON()
}
