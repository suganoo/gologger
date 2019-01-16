package gologger

import (
	"fmt"
	"bufio"
	"os"
	"testing"
	"regexp"
	"strings"
	"strconv"
	"io/ioutil"
)

var testlog = "gologger_test.log"

func init() {
	os.Remove(testlog)
}

func TestGetVersion(t *testing.T) {
	t.Log(getVersion())
	if getVersion() != version {
		t.Error("Error get version.")
	}
}

func TestSetVersion(t *testing.T) {
	nextV := "3.2.1"
	SetVersion(nextV)
	t.Log(getVersion())
	if getVersion() != version {
		t.Error("Error set version.")
	}

	SetVersion("1.0.0")
}

func TestSeparator(t *testing.T) {
	if separator != "\t" {
		t.Error("Error default separator.")
	}
}

func TestSetSeparator(t *testing.T) {
	newSp := "---"
	SetSeparator(newSp)
	if separator != newSp {
		t.Error("Error set separator.")
	}

	SetSeparator("\t")
}

func TestMuteDebug(t *testing.T) {
	glog := NewGologger(Configuration{})
	glog.MuteDebug()
	if glog.Config.ShowDebug != false {
		t.Error("Error mute debug.")
	}
}

func TestUnMuteDebug(t *testing.T) {
	glog := NewGologger(Configuration{})
	glog.UnmuteDebug()
	if glog.Config.ShowDebug != true {
		t.Error("Error unmute debug.")
	}
}

func TestTimeFormat(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	SetItemsList([]string{})
	SetTimeFormat("2006/01/02")

	r := regexp.MustCompile(`[0-9]{4}/[0-9]{2}/[0-9]{2}`)

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")

	if ! r.MatchString(dl[0]) {
		t.Error("Error time format in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
	SetTimeFormat("2006-01-02T15:04:05.000-07:00")
}

func TestLogLevel(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	glog.UnmuteDebug()
	SetItemsList([]string{KeyLogLevel})

	glog.Debug()
	glog.Info()
	glog.Warning()
	glog.Error()

	fp, err := os.Open(testlog)
	if err != nil {
		t.Error("Open error.")
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	var i = 0
	for scanner.Scan() {
		data := scanner.Text()
		t.Log(data)

		dl := strings.Split(data, "\t")
		switch i {
		case 0:
			if dl[1] != "DEBUG" {
				t.Error("Error log level DEBUG in log.")
			}
		case 1:
			if dl[1] != "INFO" {
				t.Error("Error log level INFO in log.")
			}
		case 2:
			if dl[1] != "WARNING" {
				t.Error("Error log level WARNING in log.")
			}
		case 3:
			if dl[1] != "ERROR" {
				t.Error("Error log level ERROR in log.")
			}
		}
		i++
	}
	glog.CloseFile()
	os.Remove(testlog)
}

func TestVersion(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	SetItemsList([]string{KeyVersion})

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")
	if dl[1] != "1.0.0" {
		t.Error("Error version in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestHostname(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	SetItemsList([]string{KeyHostName})
	hostname := "hogeserver"
	st.hostname = hostname

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")
	if dl[1] != hostname {
		t.Error("Error hostname in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestUsername(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	SetItemsList([]string{KeyUserName})
	username := "suganoo"
	st.username = username

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")
	if dl[1] != username {
		t.Error("Error username in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestMessage(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	SetItemsList([]string{KeyMessage})
	msg := "hogehoge"
	num := 1234
	type HogeStruct struct {
		Id   int
		Name string
	}
	var hs = HogeStruct{Id:1234,Name:"suganoo"}
	list := []int{1, 2, 3, 4}

	glog.Info(msg)
	glog.Info(num)
	glog.Info(hs)
	glog.Info(msg, num)
	glog.Info(list)

	fp, err := os.Open(testlog)
	if err != nil {
		t.Error("Open error.")
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	var i = 0
	for scanner.Scan() {
		data := scanner.Text()
		t.Log(data)

		dl := strings.Split(data, "\t")
		switch i {
		case 0:
			if dl[1] != msg {
				t.Error("Error string in log.")
			}
		case 1:
			if dl[1] != strconv.Itoa(num) {
				t.Error("Error int in log.")
			}
		case 2:
			if dl[1] != fmt.Sprintf("%v", hs) {
				t.Error("Error struct in log.")
			}
		case 3:
			if dl[1] != (msg + " " + strconv.Itoa(num)) {
				t.Error("Error multi in log.")
			}
		case 4:
			if dl[1] != fmt.Sprintf("%v", list) {
				t.Error("Error list in log.")
			}
		}
		i++
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestProcessId(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	SetItemsList([]string{KeyProcessId})

	r := regexp.MustCompile(`\d`)

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")

	if ! r.MatchString(dl[1]) {
		t.Error("Error process id in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestGoroutineId(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	SetItemsList([]string{KeyGoroutineId})

	r := regexp.MustCompile(`GrtnID:\d`)

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")

	if ! r.MatchString(dl[1]) {
		t.Error("Error goroutine id in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestFuncandFileName(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	SetItemsList([]string{KeyFuncAndFileName})

	r := regexp.MustCompile(`[gologger_test.go:\d]`)

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")

	if dl[1] != "TestFuncandFileName" {
		t.Error("Error function name in log.")
	}
	if ! r.MatchString(dl[2]) {
		t.Error("Error filename in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}
