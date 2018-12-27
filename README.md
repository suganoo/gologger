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
2018-02-21T10:07:44.277+09:00	INFO	hoge.sever	3892	fuga-user	1.0.0	this is info	main	[gologgerSample.go:18]
2018-02-21T10:07:44.277+09:00	INFO	hoge.sever	3892	fuga-user	1.0.0	msg : hogehoge	main	[gologgerSample.go:19]
2018-02-21T10:07:44.277+09:00	WARNING	hoge.sever	3892	fuga-user	1.0.0	this is warning	main	[gologgerSample.go:20]
2018-02-21T10:07:44.277+09:00	ERROR	hoge.sever	3892	fuga-user	1.0.0	this is Error	main	[gologgerSample.go:21]
2018-02-21T10:07:44.277+09:00	DEBUG	hoge.sever	3892	fuga-user	1.0.0	this is debug hogeFunc	hogeFunc	[gologgerSample.go:8]
2018-02-21T10:07:44.277+09:00	INFO	hoge.sever	3892	fuga-user	1.0.0	call hogeFunc	hogeFunc	[gologgerSample.go:9]
2018-02-21T10:07:44.277+09:00	DEBUG	hoge.sever	3892	fuga-user	1.0.0	this is debug xxx	main	[gologgerSample.go:26]
```
