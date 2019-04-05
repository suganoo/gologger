package gologger

import (
	"errors"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"strconv"
	"sync"
	"time"
)

type KeyId int
const (
	KeyTimestamp KeyId = iota
	KeyLogLevel
	KeyHostName
	KeyProcessId
	KeyGoroutineId
	KeyUserName
	KeyVersion
	KeyMessage
	KeyFunc
	KeyFileName
)

// these are keys for json format.
var knm map[KeyId]string = map[KeyId]string{}
const (
	KeyNameTimestamp   = "timestamp"
	KeyNameLogLevel    = "loglevel"
	KeyNameHostName    = "hostname"
	KeyNameProcessId   = "pid"
	KeyNameGoroutineId = "gid"
	KeyNameUserName    = "username"
	KeyNameVersion     = "version"
	KeyNameMessage     = "msg"
	KeyNameFunc        = "func"
	KeyNameFileName    = "filename"
)

type OutputFmtType int
const (
	FmtDefault   OutputFmtType = iota
	FmtJSON
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
	LogItems   []KeyId
}

type Gologger struct {
	Config     Configuration
	FormatterInterface
}

var f *(os.File)
var err error
var lock sync.Mutex

func init() {
	// these map are for json format.
	knm[KeyTimestamp]   = KeyNameTimestamp
	knm[KeyLogLevel]    = KeyNameLogLevel
	knm[KeyHostName]    = KeyNameHostName
	knm[KeyProcessId]   = KeyNameProcessId
	knm[KeyGoroutineId] = KeyNameGoroutineId
	knm[KeyUserName]    = KeyNameUserName
	knm[KeyVersion]     = KeyNameVersion
	knm[KeyMessage]     = KeyNameMessage
	knm[KeyFunc]        = KeyNameFunc
	knm[KeyFileName]    = KeyNameFileName
}

func (g Gologger) Write(bytes []byte) (int, error) {
	lock.Lock()
	defer lock.Unlock()
	msg := string(bytes)
	return f.Write(([]byte)(msg))
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

func (g *Gologger)SetItemsList(itemsList []KeyId) {
	g.Config.LogItems = itemsList
}

func (g *Gologger)SetOutputFormat(typeId OutputFmtType) {
	switch typeId {
	case FmtDefault:
		g.FormatterInterface = MarshallFunc(defaultFormat)
	case FmtJSON:
		g.FormatterInterface = MarshallFunc(jsonFormat)
	default:
		g.FormatterInterface = MarshallFunc(defaultFormat)
	}
}

// Log Format
type FormatterInterface interface {
	marshall(*Gologger, string, string) string
}

type MarshallFunc func(*Gologger, string, string) string
func (m MarshallFunc) marshall(g *Gologger, logLevel string, msg string) (logMsg string){
	return m(g, logLevel, msg)
}

func defaultFormat(g *Gologger, logLevel string, msg string) (logMsg string){

	for _, item := range g.Config.LogItems {
		switch item {
		case KeyLogLevel:
			// set log level
			logMsg = logMsg + logLevel + g.Config.Separator
		
		case KeyMessage:
			// set log message
			logMsg = logMsg + msg      + g.Config.Separator
		
		default:
			logMsg = logMsg + g.getItem(item) + g.Config.Separator
		}
	}
	return
}

func jsonFormat(g *Gologger, logLevel string, msg string) (logMsg string){

	logMap := map[string]string{}
	for _, item := range g.Config.LogItems {
		switch item {
		case KeyLogLevel:
			// set log level
			logMap[knm[item]] = logLevel
		
		case KeyMessage:
			// set log message
			logMap[knm[item]] = msg
		
		case KeyGoroutineId:
			// split goroutine id, ex.  gid:1 -> 1
			logMap[knm[item]] = strings.Split(g.getItem(item), ":")[1]
		
		default:
			logMap[knm[item]] = g.getItem(item)
		}
	}
	str, _ := json.Marshal(logMap)
	logMsg = string(str)
	return
}

func (g *Gologger)getItem(logType KeyId) (string) {

	switch logType {
	case KeyTimestamp:
		// set timestamp
		timestamp := time.Now().Format(g.Config.TimeFormat)
		return timestamp

	case KeyHostName:
		// set hostname
		//return st.getHostname()
		return g.getHostname()
	
	case KeyProcessId:
		// set process id
		pid := os.Getpid()
		return strconv.Itoa(pid)
	
	case KeyGoroutineId:
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
		return KeyNameGoroutineId + ":" + strings.Split(string(rsb)," ")[1]

	case KeyUserName:
		// set user name
		//return st.getUsername()
		return g.getUsername()
	
	case KeyVersion:
		// set version
		return g.getVersion()
	
	case KeyFunc, KeyFileName:
		// call file statement
		programCounter, filePath, fileLineNum, _ := runtime.Caller(4)
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
	
	default:
		return ""
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

	// set Formatter
	gl.FormatterInterface = MarshallFunc(defaultFormat)

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
		gl.Config.LogItems = append(gl.Config.LogItems, KeyTimestamp)
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
	logMsg := g.marshall(g, "DEBUG", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Info(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.marshall(g, "INFO", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Warning(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.marshall(g, "WARNING", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Error(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.marshall(g, "ERROR", msg[1:len(msg)-1])
	log.Println(logMsg)
}

//===== the following funcs are not recommended
func (g *Gologger)Fatal(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.marshall(g, "FATAL", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Panic(v ...interface{}) {
	msg := fmt.Sprintf("%v", v)
	logMsg := g.marshall(g, "PANIC", msg[1:len(msg)-1])
	log.Println(logMsg)
	panic("Call panic from Gologger")
}

