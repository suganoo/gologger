# suganoo/gologger
logger for Golang.

## How to use
```
go get -u github.com/suganoo/gologger
```
ex. gologgerSample.go
```Go:gologgerSample.go
package main

import (
	"github.com/suganoo/gologger"
)

var glog *gologger.Gologger

func hogeFunc() {
	glog.Debug("this is debug hogeFunc")
	glog.Info("call hogeFunc")
}

func main() {
	glog = gologger.NewGologger(gologger.Configuration{Logfile : "./testlog.log"})
	defer glog.CloseFile()

	msg := "hogehoge"
	glog.Debug("this is debug")   // default debug is muted
	glog.Info("this is info")
	glog.Info("msg : " + msg)
	glog.Warning("this is warning")
	glog.Error("this is Error")

	glog.UnmuteDebug()
	hogeFunc()

	glog.Debug("this is debug xxx")
	glog.MuteDebug()
	glog.Debug("this is debug yyy")  // this debug message is muted
}
```
```
go run gologgerSample.go
```
## Output
testlog.log
```testlog.log
2018-02-21T10:07:44.277+09:00	INFO	hoge.sever	3892	GrtnID:1	fuga-user	1.0.0	this is info	main	[gologgerSample.go:18]
2018-02-21T10:07:44.277+09:00	INFO	hoge.sever	3892	GrtnID:1	fuga-user	1.0.0	msg : hogehoge	main	[gologgerSample.go:19]
2018-02-21T10:07:44.277+09:00	WARNING	hoge.sever	3892	GrtnID:1	fuga-user	1.0.0	this is warning	main	[gologgerSample.go:20]
2018-02-21T10:07:44.277+09:00	ERROR	hoge.sever	3892	GrtnID:1	fuga-user	1.0.0	this is Error	main	[gologgerSample.go:21]
2018-02-21T10:07:44.277+09:00	DEBUG	hoge.sever	3892	GrtnID:1	fuga-user	1.0.0	this is debug hogeFunc	hogeFunc	[gologgerSample.go:8]
2018-02-21T10:07:44.277+09:00	INFO	hoge.sever	3892	GrtnID:1	fuga-user	1.0.0	call hogeFunc	hogeFunc	[gologgerSample.go:9]
2018-02-21T10:07:44.277+09:00	DEBUG	hoge.sever	3892	GrtnID:1	fuga-user	1.0.0	this is debug xxx	main	[gologgerSample.go:26]
```
## Option
### Any type of message
Any type of messages can be logged.
```
// Array
items := []string{"aaaa", "bbbb"}
glog.Info(items)
// Int
glog.Info(1000)
// Struct
type Hoge struct {
	Id int
	Name string
}
hoge := Hoge{Id: 1222, Name:"aaaa"}
glog.Info(hoge)
// Mix
glog.Info("ddddddddddd", "ooooo", 123)
```
```
// Array
2019-01-07T12:17:01.248+09:00	INFO	hoge.sever	18184	GrtnID:1	fuga-user	1.0.0	[aaaa bbbb]	main	[main.go:24]
// Int
2019-01-07T12:17:01.248+09:00	INFO	hoge.sever	18184	GrtnID:1	fuga-user	1.0.0	1000	main	[main.go:23]
// Struct
2019-01-07T12:17:01.248+09:00	INFO	hoge.sever	18184	GrtnID:1	fuga-user	1.0.0	{1222 aaaa}	main	[main.go:30]
// Mix
2019-01-07T12:17:01.248+09:00	INFO	hoge.sever	18184	GrtnID:1	fuga-user	1.0.0	ddddddddddd ooooo 123	main	[main.go:31]
```
### Time Format
It is available to change time format.
```
gologger.SetTimeFormat("2006/01/02")
```
```
2018/02/21	INFO	hoge.sever...
```
### Separator
It is also available to change separator. Default separator is tab.
```
gologger.SetSeparator("---")
```
```
2018-02-21T10:07:44.277+09:00---INFO---hoge.sever---...
```
### Items order
It is also available to change the order of log items except timestamp.
```
ex.
gologger.SetItemsList([]string{gologger.KeyMessage, gologger.KeyFuncAndFileName, gologger.KeyLogLevel, gologger.KeyPid})
glog.Info("hogehoge")
```
```
2019-01-07T12:21:33.893+09:00	hogehoge	main	[main.go:11]	INFO	10780
```
