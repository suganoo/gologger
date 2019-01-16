package gologger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"strconv"
	"time"
)

var version = "1.0.0"
var separator = "\t"
var timeFormat = "2006-01-02T15:04:05.000-07:00"

// base statement
type Statement struct {
	hostname  string
	username  string
}

// log format
type LogWriter struct {
}

// Gologger struct
type Configuration struct {
	Logfile    string
	ShowDebug  bool
}

type Gologger struct {
	Config     Configuration
}

var KeyLogLevel = "LogLevel"
var KeyHostname = "Hostname"
var KeyProcessId   = "ProcessId"
var KeyGoroutineId = "GoroutineId"
var KeyUserName = "UserName"
var KeyVersion  = "Version"
var KeyMessage  = "Message"
var KeyFuncAndFileName = "FuncAndFileName"

var st Statement
var logItems []string
var f *(os.File)
var err error

func init() {
	// set hostname, username ....
	hostname, _ := os.Hostname()
	st.hostname = hostname

	user, _  := user.Current()
	st.username = user.Username

	// log settings
	log.SetFlags(0)
	log.SetOutput(new(LogWriter))

	// log item position
	logItems = append(logItems, KeyLogLevel)
	logItems = append(logItems, KeyHostname)
	logItems = append(logItems, KeyProcessId)
	logItems = append(logItems, KeyGoroutineId)
	logItems = append(logItems, KeyUserName)
	logItems = append(logItems, KeyVersion)
	logItems = append(logItems, KeyMessage)
	logItems = append(logItems, KeyFuncAndFileName)
}

func (writer LogWriter) Write(bytes []byte) (int, error) {
	msg := string(bytes)
	timestamp := time.Now().Format(timeFormat)
	logMsg := timestamp + separator + msg

	return f.Write(([]byte)(logMsg))
}


func getVersion() string {
	return version
}

func SetVersion(vers string) {
	version = vers
}

func SetSeparator(sep string) {
	separator = sep
}

func SetTimeFormat(tf string) {
	timeFormat = tf
}

func (x Statement) getHostname() string {
	return x.hostname
}

func (x Statement) getUsername() string {
	return x.username
}

func SetItemsList(itemsList []string) {
	logItems = itemsList
}

func arrangeLog(logLevel, msg string) (logMsg string) {

	for _, item := range logItems {
		if (item == KeyLogLevel) {
			// set log level
			logMsg = logMsg + logLevel + separator
			continue
		}
		if (item == KeyMessage) {
			// set log message
			logMsg = logMsg + msg      + separator
			continue
		}
		
		logMsg = logMsg + getItem(item) + separator
	}
	return
}

func getItem(logType string) (string) {
	if (logType == KeyHostname) {
		// set hostname
		return st.getHostname()
	}
	if (logType == KeyProcessId) {
		// set process id
		pid := os.Getpid()
		return strconv.Itoa(pid)
	}
	if (logType == KeyGoroutineId) {
		// get and set goroutine id
		rsb := make([]byte, 64)
		// the content of runtime stack is like this.
		// ----------------------------
		// goroutine 1 [running]:
		// main.main()
		//     C:/.....
		runtime.Stack(rsb, false)
		// so get goroutine id
		// "goroutine 1 [running]:" --> "1"
		return "GrtnID:" + strings.Split(string(rsb)," ")[1]

	}
	if (logType == KeyUserName) {
		// set user name
		return st.getUsername()
	}
	if (logType == KeyVersion) {
		// set version
		return getVersion()
	}
	if (logType == KeyFuncAndFileName) {
		// call file statement
		programCounter, filePath, fileLineNum, _ := runtime.Caller(3)
		filePathArry := strings.Split(fmt.Sprintf("%v",filePath), "/")
	
		// set called function name
		fn := runtime.FuncForPC(programCounter)
		fnNameArry := strings.Split(fn.Name(), ".")
	
		// set function and filename with line number
		return fnNameArry[1]     + separator + "[" + filePathArry[len(filePathArry) - 1] + ":" + strconv.Itoa(fileLineNum) + "]"
	}
	return ""
}

func NewGologger(conf Configuration) (*Gologger) {
	gl := &Gologger{
		Config: conf,
	}

	f, err = os.OpenFile(gl.Config.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		errors.New("Error opening log file!" + err.Error())
	}

	return gl
}

func (g *Gologger)MuteDebug() {
	g.Config.ShowDebug = false
}

func (g *Gologger)UnmuteDebug() {
	g.Config.ShowDebug = true
}

func (g *Gologger)CloseFile() {
	f.Close()
}

func (g *Gologger)Debug(v ...interface{}) {
        if ! g.Config.ShowDebug { return }

	msg := fmt.Sprintf("%v", v)
	logMsg := arrangeLog("DEBUG", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Info(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := arrangeLog("INFO", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Warning(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := arrangeLog("WARNING", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Error(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := arrangeLog("ERROR", msg[1:len(msg)-1])
	log.Println(logMsg)
}

//===== the following funcs are not recommended
func (g *Gologger)Fatal(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := arrangeLog("FATAL", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Panic(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := arrangeLog("PANIC", msg[1:len(msg)-1])
	log.Println(logMsg)
	panic("Call panic from Gologger")
}

