package conf

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
	"github.com/qiusanshiye/webDialTestTool/logs"
)

type TAPIConf struct {
	Name       string `ini:"-"`
	Addr       string `ini:"addr"`
	Interval   int    `ini:"interval"` // unit: s
	Uri        string `ini:"uri"`
	Query      string `ini:"query"`
	Body       string `ini:"body"`
	CostLine   int    `ini:"costline"` // unit: ms
	DDToken    string `ini:"ddtoken"`
	ReqHeaders map[string]string
}

var (
	GAPIConfs = map[string]TAPIConf{}
	GLogConf  logs.TLoggerConf
)

const (
	cfgAPIBaseName = "API"
	cfgLogName     = "LOGS"
)

func (cf TAPIConf) String() string {
	return fmt.Sprintf(
		"TAPIConf{addr:%s, inter:%d, uri:%s, query=%s, costline:%d, heads=%+v}",
		cf.Addr, cf.Interval, cf.Uri, cf.Query, cf.CostLine, cf.ReqHeaders)
}

func InitConfig() {
	iniPath := flag.String("c", "config.ini", "-c conf-path")
	flag.Parse()

	fp, err := ini.LoadSources(
		ini.LoadOptions{UnescapeValueDoubleQuotes: true},
		*iniPath)
	if err != nil {
		checkError(err)
	}

	if err := loadLogConf(fp); err != nil {
		checkError(err)
	}
	if err := loadAPIConf(fp); err != nil {
		checkError(err)
	}
}

func loadLogConf(fp *ini.File) error {
	return fp.Section(cfgLogName).MapTo(&GLogConf)
}
func loadAPIConf(fp *ini.File) error {
	defaultSec := fp.Section(cfgAPIBaseName)
	secs := defaultSec.ChildSections()
	size := len(secs)
	if size == 0 {
		return fmt.Errorf("no api configs found")
	}

	for _, sec := range secs {
		name := sec.Name()
		GAPIConfs[name] = parseSectionToStruct(sec, defaultSec)
	}
	return nil
}

func parseSectionToStruct(sec, defaultSec *ini.Section) TAPIConf {
	defaultKeyHash := defaultSec.KeysHash()
	keyHash := sec.KeysHash()

	cf := TAPIConf{}

	var exists bool
	if cf.Addr, exists = keyHash["addr"]; !exists {
		cf.Addr = defaultKeyHash["addr"]
	}
	if cf.Uri, exists = keyHash["uri"]; !exists {
		cf.Uri = defaultKeyHash["uri"]
	}
	if cf.Query, exists = keyHash["query"]; !exists {
		cf.Query = defaultKeyHash["query"]
	}
	if cf.Body, exists = keyHash["body"]; !exists {
		cf.Body = defaultKeyHash["body"]
	}
	if cf.DDToken, exists = keyHash["ddtoken"]; !exists {
		cf.DDToken = defaultKeyHash["ddtoken"]
	}

	if _, exists = keyHash["interval"]; exists {
		cf.Interval, _ = strconv.Atoi(keyHash["interval"])
	} else {
		cf.Interval, _ = strconv.Atoi(defaultKeyHash["interval"])
	}
	if _, exists = keyHash["costline"]; exists {
		cf.CostLine, _ = strconv.Atoi(keyHash["costline"])
	} else {
		cf.CostLine, _ = strconv.Atoi(defaultKeyHash["costline"])
	}

	cf.ReqHeaders = make(map[string]string)
	for k, v := range defaultKeyHash {
		if strings.HasPrefix(k, "h_") {
			headName := strings.SplitN(k, "_", 2)[1]
			cf.ReqHeaders[headName] = v
		}
	}
	for k, v := range keyHash {
		if strings.HasPrefix(k, "h_") {
			headName := strings.SplitN(k, "_", 2)[1]
			cf.ReqHeaders[headName] = v
		}
	}

	return cf
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
