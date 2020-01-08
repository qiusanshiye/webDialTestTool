package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/qiusanshiye/webDialTestTool/conf"
	"github.com/qiusanshiye/webDialTestTool/logs"
)

var wg = new(sync.WaitGroup)

func main() {
	conf.InitConfig()

	logs.Init(&conf.GLogConf)
	defer logs.Uninit()

	dialTestRun()
}

func dialTestRun() {
	nc := make(chan bool)
	for k, v := range conf.GAPIConfs {
		logs.Info("dialTest: %s: %s", k, v)

		wg.Add(1)
		go apiDialTestStart(k, v, nc)
	}

	wg.Wait()
	logs.Info("test exit. bye-bye")
}

func apiDialTestStart(title string, apiConf conf.TAPIConf, nc chan bool) {
	logs.Info("title:%s api:%s dial-test-start. interval %d",
		title, apiConf.Uri, apiConf.Interval)
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(apiConf.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			logs.Info("title:%s api:%s test begin", title, apiConf.Uri)
			apiDialTestImp(title, apiConf)
		case s, ok := <-nc:
			if !ok {
				logs.Warn("signal chan closed, exit")
				break
			}
			logs.Warn("recv signal %v, exit", s)
			break
		}
	}
}

func apiDialTestImp(title string, apiConf conf.TAPIConf) {
	url := fmt.Sprintf("%s%s", apiConf.Addr, apiConf.Uri)
	if len(apiConf.Query) > 0 {
		url = fmt.Sprintf("%s%s?%s", apiConf.Addr, apiConf.Uri, apiConf.Query)
	}

	var postBody []byte = nil
	if len(apiConf.Body) > 0 {
		postBody = []byte(apiConf.Body)
	}

	beg := time.Now().UnixNano() / 1000000
	code, resp, err := httpRequest(url, apiConf.ReqHeaders, postBody)
	end := time.Now().UnixNano() / 1000000
	cost := int(end - beg)
	if err != nil {
		logs.Error("title:%s api:%s dial-test failed! err=%s", title, apiConf.Uri, err)
		failAlarm(apiConf.DDToken, title, apiConf.Uri, string(resp), err, cost)
		return
	}

	if code != 200 {
		logs.Error("title:%s api:%s dial-test failed! code=%d", title, apiConf.Uri, code)
		err := fmt.Errorf("resp status:%d, non 200", code)
		failAlarm(apiConf.DDToken, title, apiConf.Uri, string(resp), err, cost)
	}

	if cost > apiConf.CostLine {
		logs.Error("title:%s api:%s test cost time: %d", title, apiConf.Uri, cost)
		timeoutAlarm(apiConf.DDToken, title, apiConf.Uri, string(resp), nil, cost)
	}
	logs.Info("title:%s api:%s dial-test done. cost=%d, resp=%s", title, apiConf.Uri, cost, resp)
}

func failAlarm(ddtoken, name, uri, resp string, err error, cost int) error {
	title := fmt.Sprintf("拨测:%s 失败", uri)

	content := fmt.Sprintf("## %s", title)
	content += fmt.Sprintf("\n- --------------------------")
	content += fmt.Sprintf("\n- time=%s", time.Now().Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("\n- name=%s", name)
	content += fmt.Sprintf("\n- resp=%s", resp)
	if err != nil {
		content += fmt.Sprintf("\n- err=%v", err)
	}
	content += fmt.Sprintf("\n- cost=%dms", cost)
	content += fmt.Sprintf("\n- --------------------------")

	msg, err := alarmNotify(title, content, ddtoken)
	if err != nil {
		logs.Error("alarm: title:%s failed, err=%s", title, err)
		return err
	}
	logs.Debug("alarm success: title:%s result=%s", title, msg)
	return nil
}

func timeoutAlarm(ddtoken, name, uri, resp string, err error, cost int) error {
	title := fmt.Sprintf("拨测:%s 超时", uri)

	content := fmt.Sprintf("## %s", title)
	content += fmt.Sprintf("\n- --------------------------")
	content += fmt.Sprintf("\n- time=%s", time.Now().Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("\n- name=%s", name)
	content += fmt.Sprintf("\n- cost=%dms", cost)
	content += fmt.Sprintf("\n- resp=%s", resp)
	if err != nil {
		content += fmt.Sprintf("\n- err=%v", err)
	}
	content += fmt.Sprintf("\n- --------------------------")

	msg, err := alarmNotify(title, content, ddtoken)
	if err != nil {
		logs.Error("alarm: title:%s failed, err=%s", title, err)
		return err
	}
	logs.Debug("alarm success: title:%s result=%s", title, msg)
	return nil
}
