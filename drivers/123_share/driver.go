/*
 * drivers/123_share/driver.go
 */
package _123Share

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/json-iterator/go"

	"github.com/alist-org/alist/v3/internal/driver/base"
	"github.com/alist-org/alist/v3/pkg/model"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/pkg/log"
)

func (d *Pan123Share) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	if f, ok := file.(File); ok {
		data := base.Json{
			"shareKey":  d.ShareKey,
			"SharePwd":  d.SharePwd,
			"etag":      f.Etag,
			"fileId":    f.FileId,
			"s3keyFlag": f.S3KeyFlag,
			"size":      f.Size,
		}
		var headers map[string]string
		if !utils.IsLocalIPAddr(args.IP) {
			headers = map[string]string{
				"X-Forwarded-For": args.IP,
			}
		}
		resp, err := d.request(DownloadInfo, http.MethodPost, func(req *resty.Request) {
			req.SetBody(data).
				SetHeaders(headers).
				SetHeader("User-Agent", "123pan/v2.4.0(Android_10.0;Xiaomi)").
				SetHeader("platform", "android").
				SetHeader("app-version", "61").
				SetHeader("x-app-version", "2.4.0")
		}, nil)
		if err != nil {
			return nil, err
		}
		downloadUrl := utils.Json.Get(resp, "data", "DownloadURL").ToString()
		u, err := url.Parse(downloadUrl)
		if err != nil {
			return nil, err
		}
		nu := u.Query().Get("params")
		if nu != "" {
			du, _ := base64.StdEncoding.DecodeString(nu)
			u, err = url.Parse(string(du))
			if err != nil {
				return nil, err
			}
			q := u.Query()
			q.Set("auto_redirect", "0")
			u.RawQuery = q.Encode()
		}
		u_ := u.String()
		log.Debug("download url: ", u_)
		res, err := base.NoRedirectClient.R().
			SetHeader("Referer", "https://www.123pan.com/").
			SetHeader("User-Agent", "123pan/v2.4.0(Android_10.0;Xiaomi)").
			SetHeader("platform", "android").
			SetHeader("app-version", "61").
			SetHeader("x-app-version", "2.4.0").
			Get(u_)
		if err != nil {
			return nil, err
		}
		link := model.Link{URL: u_}
		if res.StatusCode() == 302 {
			link.URL = res.Header().Get("location")
		} else if res.StatusCode() < 300 {
			link.URL = utils.Json.Get(res.Body(), "data", "redirect_url").ToString()
		}
		link.Header = http.Header{
			"Referer": []string{"https://www.123pan.com/"},
		}
		return &link, nil
	}
	return nil, fmt.Errorf("can't convert obj")
}
