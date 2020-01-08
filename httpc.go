package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/qiusanshiye/webDialTestTool/logs"
)

var DefaultTransport http.RoundTripper = &http.Transport{
	// Proxy: ProxyFromEnvironment, //代理使用
	DialContext: (&net.Dialer{
		Timeout:   10 * time.Second, //连接超时时间
		KeepAlive: 30 * time.Second, //连接保持超时时间
		DualStack: true,             //
	}).DialContext,
	MaxIdleConns:          500, //client对与所有host最大空闲连接数总和
	MaxIdleConnsPerHost:   10,
	ResponseHeaderTimeout: 15 * time.Second,
	IdleConnTimeout:       90 * time.Second, //空闲连接在连接池中的超时时间
	TLSHandshakeTimeout:   10 * time.Second, //TLS安全连接握手超时时间
	ExpectContinueTimeout: 3 * time.Second,  //发送完请求到接收到响应头的超时时间
}

func httpRequest(url string, headers map[string]string, body []byte) (int, []byte, error) {
	var data *bytes.Reader = nil
	var method string = "GET"

	if body != nil {
		method = "POST"
		data = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, data)
	if err != nil {
		logs.Error("http new request failed. err=%s", err)
		return 0, nil, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	httpc := &http.Client{
		Transport: DefaultTransport,
	}
	resp, err := httpc.Do(req)
	if err != nil {
		logs.Error("http do failed. err=%s", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp.StatusCode, nil, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read resp body failed. err=%s", err)
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, respBody, nil
}
