/*
 * drivers/123_share/util.go
 */
package _123Share

import (
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/alist-org/alist/v3/internal/driver/base"
	"github.com/alist-org/alist/v3/pkg/utils"
)

func (d *Pan123Share) request(url string, method string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	req := base.RestyClient.R()
	req.SetHeaders(map[string]string{
		"origin":          "https://www.123pan.com",
		"referer":         "https://www.123pan.com/",
		"authorization":   "Bearer " + d.AccessToken,
		"user-agent":      "123pan/v2.4.0(Android_10.0;Xiaomi)",
		"platform":        "android",
		"app-version":     "61",
		"x-app-version":   "2.4.0",
	})
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	res, err := req.Execute(method, GetApi(url))
	if err != nil {
		return nil, err
	}
	body := res.Body()
	code := utils.Json.Get(body, "code").ToInt()
	if code != 0 {
		return nil, errors.New(jsoniter.Get(body, "message").ToString())
	}
	return body, nil
}
