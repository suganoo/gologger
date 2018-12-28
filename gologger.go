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

type PositionItems struct {
	PosLogLevel   int
	PosHostname   int
	PosPid        int
	PosUserName   int
	PosVersion    int
	PosMessage    int
	PosFuncName   int
	PosFileName   int
}

var st Statement
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

func arrangeLog(logLevel, msg string) (logMsg string) {

	// set log level
	logMsg = logLevel + separator

	// set hostname
	logMsg = logMsg + st.getHostname()  + separator

	// set process id
	pid := os.Getpid()
	logMsg = logMsg + strconv.Itoa(pid) + separator

	// set user name
	logMsg = logMsg + st.getUsername()  + separator

	// set version
	logMsg = logMsg + getVersion()      + separator

	// set log message
	logMsg = logMsg + msg               + separator

	// call file statement
	programCounter, filePath, fileLineNum, _ := runtime.Caller(2)
	filePathArry := strings.Split(fmt.Sprintf("%v",filePath), "/")

	// set called function name
	fn := runtime.FuncForPC(programCounter)
	fnNameArry := strings.Split(fn.Name(), ".")

	logMsg = logMsg + fnNameArry[1]     + separator

	// set filename
	logMsg = logMsg + "[" + filePathArry[len(filePathArry) - 1] + ":" + strconv.Itoa(fileLineNum) + "]"

	return
}

func NewGologger(conf Configuration) (*Gologger) {
	gl := &Gologger{
		Config: conf,
	}
	fmt.Println(gl.Config.ShowDebug)

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

	msg := fmt.Sprint(v)
	logMsg := arrangeLog("DEBUG", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Info(v ...interface{}) {
	msg := fmt.Sprint(v)
	logMsg := arrangeLog("INFO", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Warning(v ...interface{}) {
	msg := fmt.Sprint(v)
	logMsg := arrangeLog("WARNING", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Error(v ...interface{}) {
	msg := fmt.Sprint(v)
	logMsg := arrangeLog("ERROR", msg[1:len(msg)-1])
	log.Println(logMsg)
}

//===== the following funcs are not recommended
func (g *Gologger)Fatal(v ...interface{}) {
	msg := fmt.Sprint(v)
	logMsg := arrangeLog("FATAL", msg[1:len(msg)-1])
	log.Println(logMsg)
}

func (g *Gologger)Panic(v ...interface{}) {
	msg := fmt.Sprint(v)
	logMsg := arrangeLog("PANIC", msg[1:len(msg)-1])
	log.Println(logMsg)
	panic("Call panic from Gologger")
}

