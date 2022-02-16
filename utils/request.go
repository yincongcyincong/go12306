package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/tools/12306/module"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var client *http.Client

func GetClient() *http.Client {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 1 * time.Minute,
				}).DialContext,
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   20,
				IdleConnTimeout:       10 * time.Second,
			},
		}
	}

	return client

}

func GetCdnClient(cdn string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				addr = cdn + ":443"
				return (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 1 * time.Minute,
				}).DialContext(ctx, network, addr)
			},
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   20,
			IdleConnTimeout:       10 * time.Second,
		},
	}
}

func Request(data string, cookieStr, url string, res interface{}, headers map[string]string) error {

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
	req.Header.Set("Host", "kyfw.12306.cn")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := GetClient().Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, res)
	if err != nil {
		seelog.Errorf("json unmarshal fail: %v, %v, %v", err, string(respBody), url)
		return err
	}

	// 添加cookie
	setCookies := resp.Header.Values("Set-Cookie")
	AddCookieStr(setCookies)

	seelog.Tracef("url: %v, param: %v, response: %v", url, data, string(respBody))

	return nil
}

func RequestGet(cookieStr, url string, res interface{}, headers map[string]string) error {

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
	req.Header.Set("Host", "kyfw.12306.cn")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := GetClient().Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, res)
	if err != nil {
		return err
	}

	// 添加cookie
	setCookies := resp.Header.Values("Set-Cookie")
	AddCookieStr(setCookies)

	seelog.Tracef("url: %v, response: %v", url, string(respBody))

	return nil
}

func RequestGetWithoutJson(cookieStr, url string, headers map[string]string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
	req.Header.Set("Host", "kyfw.12306.cn")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := GetClient().Do(req)
	if err != nil {
		return []byte{}, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	// 添加cookie
	setCookies := resp.Header.Values("Set-Cookie")
	AddCookieStr(setCookies)

	seelog.Tracef("url: %v, response: %v", url, string(respBody))

	return respBody, nil
}

func RequestGetWithCDN(cookieStr, url string, res interface{}, headers map[string]string, cdn string) error {

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
	req.Header.Set("Host", "kyfw.12306.cn")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://kyfw.12306.cn")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := GetCdnClient(cdn).Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if res != nil {
		err = json.Unmarshal(respBody, res)
		if err != nil {
			return err
		}
	}

	// 添加cookie
	setCookies := resp.Header.Values("Set-Cookie")
	AddCookieStr(setCookies)

	seelog.Tracef("url: %v, response: %v, cdn: %s", url, string(respBody[:500]), cdn)

	return nil
}

func EncodeParam(r *http.Request, param interface{}) error {
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		seelog.Error(err)
		return err
	}

	err = json.Unmarshal(respBody, param)
	if err != nil {
		seelog.Error(err)
		return err
	}

	return nil
}

func HTTPSuccResp(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	resp := &module.CommonResp{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	respJson, _ := json.Marshal(resp)
	fmt.Fprint(w, string(respJson))
}

func HTTPFailResp(w http.ResponseWriter, statusCode, code int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	resp := &module.CommonResp{
		Code:    code,
		Message: message,
		Data:    data,
	}
	respJson, _ := json.Marshal(resp)
	fmt.Fprint(w, string(respJson))
}
