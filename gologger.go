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

type Statement struct {
	hostname  string
	username  string
}

type Configuration struct {
	Logfile    string
	ShowDebug  bool
	st         Statement
	Version    string
	Separator  string
	TimeFormat string
	LogItems   []string
}

type Gologger struct {
	Config     Configuration
}

var KeyLogLevel = "LogLevel"
var KeyHostName = "HostName"
var KeyProcessId   = "ProcessId"
var KeyGoroutineId = "GoroutineId"
var KeyUserName = "UserName"
var KeyVersion  = "Version"
var KeyMessage  = "Message"
var KeyFunc     = "Func"
var KeyFileName = "FileName"

var f *(os.File)
var err error

func (g Gologger) Write(bytes []byte) (int, error) {
	msg := string(bytes)
	timestamp := time.Now().Format(g.Config.TimeFormat)
	logMsg := timestamp + g.Config.Separator + msg

	return f.Write(([]byte)(logMsg))
}

func (g *Gologger)getHostname() string {
	return g.Config.st.hostname
}

func (g *Gologger)getUsername() string {
	return g.Config.st.username
}

func (g *Gologger)getVersion() string {
	return g.Config.Version
}

func (g *Gologger)SetVersion(vers string) {
	g.Config.Version = vers
}

func (g *Gologger)SetSeparator(sep string) {
	g.Config.Separator = sep
}

func (g *Gologger)SetTimeFormat(tf string) {
	g.Config.TimeFormat = tf
}

func (g *Gologger)SetItemsList(itemsList []string) {
	g.Config.LogItems = itemsList
}

func (g *Gologger)arrangeLog(logLevel, msg string) (logMsg string) {

	for _, item := range g.Config.LogItems {
		if (item == KeyLogLevel) {
			// set log level
			logMsg = logMsg + logLevel + g.Config.Separator
			continue
		}
		if (item == KeyMessage) {
			// set log message
			logMsg = logMsg + msg      + g.Config.Separator
			continue
		}
		
		logMsg = logMsg + g.getItem(item) + g.Config.Separator
	}
	return
}

func (g *Gologger)getItem(logType string) (string) {
	if (logType == KeyHostName) {
		// set hostname
		//return st.getHostname()
		return g.getHostname()
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
		//return st.getUsername()
		return g.getUsername()
	}
	if (logType == KeyVersion) {
		// set version
		return g.getVersion()
	}
	if (logType == KeyFunc) || (logType == KeyFileName){
		// call file statement
		programCounter, filePath, fileLineNum, _ := runtime.Caller(3)
		filePathArry := strings.Split(fmt.Sprintf("%v",filePath), "/")
	
		if (logType == KeyFunc){
			// set called function name
			fn := runtime.FuncForPC(programCounter)
			fnNameArry := strings.Split(fn.Name(), ".")

			return fnNameArry[1]
		}
		if (logType == KeyFileName){
			// set filename with line number
			return "[" + filePathArry[len(filePathArry) - 1] + ":" + strconv.Itoa(fileLineNum) + "]"
		}
	}
	return ""
}

func NewGologger(conf Configuration) (*Gologger) {
	gl := &Gologger{
		Config: conf,
	}

	if gl.Config.Logfile == "" {
		f = os.Stdout
	} else {
		f, err = os.OpenFile(gl.Config.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
		if err != nil {
			errors.New("Error opening log file!" + err.Error())
			os.Exit(1)
		}
	}

	// set hostname, username ....
	hostname, _ := os.Hostname()
	gl.Config.st.hostname = hostname

	user, _  := user.Current()
	gl.Config.st.username = user.Username

	// version
	if gl.Config.Version == "" {
		gl.Config.Version = "1.0.0"
	}

	// separator
	if gl.Config.Separator == "" {
		gl.Config.Separator = "\t"
	}

	// time format
	if gl.Config.TimeFormat == "" {
		gl.Config.TimeFormat = "2006-01-02T15:04:05.000-07:00"
	}

	// set log item
	if gl.Config.LogItems == nil {
		gl.Config.LogItems = append(gl.Config.LogItems, KeyLogLevel)
		gl.Config.LogItems = append(gl.Config.LogItems, KeyHostName)
		gl.Config.LogItems = append(gl.Config.LogItems, KeyProcessId)
		gl.Config.LogItems = append(gl.Config.LogItems, KeyGoroutineId)
		gl.Config.LogItems = append(gl.Config.LogItems, KeyUserName)
		gl.Config.LogItems = append(gl.Config.LogItems, KeyVersion)
		gl.Config.LogItems = append(gl.Config.LogItems, KeyMessage)
		gl.Config.LogItems = append(gl.Config.LogItems, KeyFunc)
		gl.Config.LogItems = append(gl.Config.LogItems, KeyFileName)
	}

	// log settings
	log.SetFlags(0)
	log.SetOutput(gl)

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
	logMsg := g.arrangeLog("DEBUG", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Info(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.arrangeLog("INFO", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Warning(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.arrangeLog("WARNING", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Error(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.arrangeLog("ERROR", msg[1:len(msg)-1])
	log.Println(logMsg)
}

//===== the following funcs are not recommended
func (g *Gologger)Fatal(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.arrangeLog("FATAL", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Panic(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.arrangeLog("PANIC", msg[1:len(msg)-1])
	log.Println(logMsg)
	panic("Call panic from Gologger")
}

