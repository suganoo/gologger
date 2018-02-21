# suganoo/gologger
logger for Golang.

## How to use
```
go get -u github.com/suganoo/gologger
```
```Go:gologgerSample.go
package main

import (
	"github.com/suganoo/gologger"
)

func main() {
	gologger.SetLogfile("./testlog.log")
	defer gologger.CloseFile()

	msg := "hogehoge"
	gologger.Info("this is info")
	gologger.Info("msg : " + msg)
	gologger.Warning("this is warning")
	gologger.Error("this is Error")
}
```
## Output
```testlog.log
2018-02-21T10:07:44.277+09:00	INFO	hoge.sever	3892	fuga	1.0.0	this is info	main	[gologgerSample.go:12]
2018-02-21T10:07:44.277+09:00	INFO	hoge.sever	3892	fuga	1.0.0	msg : hogehoge	main	[gologgerSample.go:13]
2018-02-21T10:07:44.277+09:00	WARNING	hoge.sever	3892	fuga	1.0.0	this is warning	main	[gologgerSample.go:14]
2018-02-21T10:07:44.277+09:00	ERROR	hoge.sever	3892	fuga	1.0.0	this is Error	main	[gologgerSample.go:15]
```
