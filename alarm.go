package main

import (
	"encoding/json"
	"fmt"
)

const (
	urlPrex = "https://oapi.dingtalk.com/robot/send?access_token="
)

type DingAlarmData struct {
	Msgtype string           `json:"msgtype"`
	Content DingAlarmContent `json:"markdown"`
}

type DingAlarmContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func alarmNotify(title, content, ddToken string) ([]byte, error) {
	url := urlPrex + ddToken

	data := DingAlarmData{
		Msgtype: "markdown",
		Content: DingAlarmContent{
			Title: title,
			Text:  content,
		},
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	heads := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 6.1; Win64; x64) Chrome/62.0.3202.94",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Content-Type":    "application/json;charset=utf-8",
	}

	code, resp, err := httpRequest(url, heads, body)
	if code != 200 {
		return nil, fmt.Errorf("resp status:%d, non 200", code)
	}
	return resp, err
}
