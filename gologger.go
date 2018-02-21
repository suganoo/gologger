package gologger

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"strconv"
	"time"
)

// base statement
type Statement struct {
	hostname  string
	username  string
}

func (x Statement) getHostname() string {
	return x.hostname
}

func (x Statement) getUsername() string {
	return x.username
}

var st Statement
var keyLevel       string = "level"
var keyMessage     string = "msg"
var keyVersion     string = "version"
var keyHostName    string = "hostname"
var keyUserName    string = "user"
var keyPid         string = "pid"
var keyFileName    string = "filename"
var keyFileLineNum string = "filelinenum"
var keyFuncName    string = "funcname"

const version = "1.0.0"
func getVersion() string {
	return version
}

// separatorとして下記の文字列に意味は無い。確率的に低いと推測して下記の文字列にした。
// The characters of separator have no meaning, just used low probability of character conbination for separator.
const separator = ":x%@:"
const separator_inner = "=<$@%="

// log format
type logWriter struct {
}

var f *(os.File)
var err error
func SetLogfile(path string) {
	f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		Error("Error opening log file!" + err.Error())
	}
}

func CloseFile() {
	f.Close()
}

func (writer logWriter) Write(bytes []byte) (int, error) {

	// separte message
	msg := string(bytes)
	msg  = strings.TrimRight(msg, "\n")
	msgArry := strings.Split(msg, separator)
	paramsMap := make(map[string]string)
	for _, item := range msgArry {
		itemArry := strings.Split(item, separator_inner)
		paramsMap[itemArry[0]] = itemArry[1]
	}

	// assemble log message
	var logMsg string
	outputSeparator := "\t"
	logMsg =          paramsMap[keyLevel]     + outputSeparator
	logMsg = logMsg + paramsMap[keyHostName]  + outputSeparator
	logMsg = logMsg + paramsMap[keyPid]       + outputSeparator
	logMsg = logMsg + paramsMap[keyUserName]  + outputSeparator
	logMsg = logMsg + paramsMap[keyVersion]   + outputSeparator
	logMsg = logMsg + paramsMap[keyMessage]   + outputSeparator
	logMsg = logMsg + paramsMap[keyFuncName]  + outputSeparator
	logMsg = logMsg + "[" + paramsMap[keyFileName] + ":" + paramsMap[keyFileLineNum] + "]"
	timestamp := time.Now().Format("2006-01-02T15:04:05.000-07:00")
	logMsg = timestamp + outputSeparator + logMsg

	return f.Write(([]byte)(logMsg + "\n"))
}

func init() {
	// set hostname, username ....
	hostname, _ := os.Hostname()
	st.hostname = hostname

	user, _  := user.Current()
	st.username = user.Username

	// log settings
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

func getParameters() (params string) {

	// set hostname
	params = params + keyHostName     + separator_inner + st.getHostname()  + separator

	// set user name
	params = params + keyUserName     + separator_inner + st.getUsername()  + separator

	// set version
	params = params + keyVersion      + separator_inner + getVersion()      + separator

	// set process id
	pid := os.Getpid()
	params = params + keyPid          + separator_inner + strconv.Itoa(pid) + separator

	// can not get thread name....

	// call file statement
	programCounter, filePath, fileLineNum, _ := runtime.Caller(2)
	filePathArry := strings.Split(fmt.Sprintf("%v",filePath), "/")

	// set filename
	params = params + keyFileName     + separator_inner + filePathArry[len(filePathArry) - 1] + separator
	// set file line number
	params = params + keyFileLineNum  + separator_inner + strconv.Itoa(fileLineNum) + separator

	// set called function name
	fn := runtime.FuncForPC(programCounter)
	fnNameArry := strings.Split(fn.Name(), ".")
	funcName := fnNameArry[1]
	params = params + keyFuncName     + separator_inner + funcName + separator
	return
}

func Info(msg string) {
	params := getParameters()
	msg = keyLevel + separator_inner + "INFO"    + separator + params + keyMessage + separator_inner + msg
	log.Println(msg)
}
func Warning(msg string) {
	params := getParameters()
	msg = keyLevel + separator_inner + "WARNING" + separator + params + keyMessage + separator_inner + msg
	log.Println(msg)
}
func Error(msg string) {
	params := getParameters()
	msg = keyLevel + separator_inner + "ERROR"   + separator + params + keyMessage + separator_inner + msg
	log.Println(msg)
}
func Fatal(msg string) {
	params := getParameters()
	msg = keyLevel + separator_inner + "ERROR"   + separator + params + keyMessage + separator_inner + msg
	log.Fatal(msg)
}
func Panic(msg string) {
	params := getParameters()
	msg = keyLevel + separator_inner + "ERROR"   + separator + params + keyMessage + separator_inner + msg
	log.Println(msg)
	panic("log panic")
}

