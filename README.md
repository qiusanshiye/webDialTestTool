# webDialTestTool
web server自动化拨测工具，支持告警与可自定义请求头

# 使用说明 
```
go get -u github.com/qiusanshiye/webDialTestTool
${GOPATH}/webDialTestTool -c <path-to>/webDialTestTool.ini
```

# 项目依赖
```
github.com/alecthomas/log4go
github.com/go-ini/ini
```

# 配置说明
```
$ cat <path-to>/webDialTestTool.ini
[LOGS]
# level: {"0: FNST", "1: FINE", "2: DEBG", "3: TRAC", "4: INFO", "5: WARN", "6: EROR", "7: CRIT"}
level=1
file=./logs/webDialTestTool.log

[API]
# addr: 默认服务器地址
addr=http://blog.5941188.com
# interval: 默认拨测间隔
interval=180
# costline: 请求耗时阀值，请求耗时超过该值时会触发钉钉告警
costline=5000
# ddtoken: 告警token, 用于请求失败、响应非200、超时告警
ddtoken=d5132d78387e0ba5c352ec7207f3a6db2c234dabf73c13604fc1f4a7d0bea21a
# h_: 自定义头，固定以  "h_" 开头的配置项做为请求头处理
h_content_type=application/json;charset=utf-8
h_connection=close

# 以下是具体api配置，均做为  [API] 的子配置存在, 即必须以 "API." 开头
# API中的所有项均可以在子配置中重新定义，否则取API中的配置值 
[API.user_login]
uri=/user/login
body={"username":"qiusanshiye","password":"123456"}

[API.user_status]
uri=/user/status
query = username=abc&passwd=123456
```
