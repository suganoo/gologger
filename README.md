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

func hogeFunc() {
	gologger.Debug("this is debug hogeFunc")
	gologger.Info("call hogeFunc")
}

func main() {
	gologger.SetLogfile("./testlog.log")
	defer gologger.CloseFile()

	msg := "hogehoge"
	gologger.Debug("this is debug")   // default debug is muted
	gologger.Info("this is info")
	gologger.Info("msg : " + msg)
	gologger.Warning("this is warning")
	gologger.Error("this is Error")

	gologger.UnmuteDebug()
	hogeFunc()

	gologger.Debug("this is debug xxx")
	gologger.MuteDebug()
	gologger.Debug("this is debug yyy")  // this debug message is muted
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
