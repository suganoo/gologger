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
var glog *Gologger

func init() {
	os.Remove(testlog)
	glog = NewGologger(Configuration{Logfile : testlog})
}

func TestGetVersion(t *testing.T) {
	t.Log(glog.getVersion())
	//if strings.Split(glog.getVersion(), "\t")[0] != "1.0.0" {
	if glog.getVersion() != "1.0.0" {
		t.Error("Error get version.")
	}
}

func TestSetVersion(t *testing.T) {
	nextV := "3.2.1"
	glog.SetVersion(nextV)
	t.Log(glog.getVersion())
	if glog.getVersion() != nextV {
		t.Error("Error set version.")
	}

	glog.SetVersion("1.0.0")
}

func TestSeparator(t *testing.T) {
	if glog.Config.Separator != "\t" {
		t.Error("Error default separator.")
	}
}

func TestSetSeparator(t *testing.T) {
	newSp := "---"
	glog.SetSeparator(newSp)
	if glog.Config.Separator != newSp {
		t.Error("Error set separator.")
	}

	glog.SetSeparator("\t")
}

func TestMuteDebug(t *testing.T) {
	glog.MuteDebug()
	if glog.Config.ShowDebug != false {
		t.Error("Error mute debug.")
	}
}

func TestUnMuteDebug(t *testing.T) {
	glog.UnmuteDebug()
	if glog.Config.ShowDebug != true {
		t.Error("Error unmute debug.")
	}
}

func TestTimeFormat(t *testing.T) {
	glog.CloseFile()
	os.Remove(testlog)

	glog = NewGologger(Configuration{Logfile: testlog})
	glog.SetItemsList([]KeyId{KeyTimestamp})
	glog.SetTimeFormat("2006/01/02")

	r := regexp.MustCompile(`[0-9]{4}/[0-9]{2}/[0-9]{2}`)

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))

	if ! r.MatchString(string(data)) {
		t.Error("Error time format in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
	glog.SetTimeFormat("2006-01-02T15:04:05.000-07:00")
}

func TestLogLevel(t *testing.T) {
	os.Remove(testlog)

	glog := NewGologger(Configuration{Logfile: testlog})
	glog.UnmuteDebug()
	glog.SetItemsList([]KeyId{KeyLogLevel})

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
			if dl[0] != "DEBUG" {
				t.Error("Error log level DEBUG in log.")
			}
		case 1:
			if dl[0] != "INFO" {
				t.Error("Error log level INFO in log.")
			}
		case 2:
			if dl[0] != "WARNING" {
				t.Error("Error log level WARNING in log.")
			}
		case 3:
			if dl[0] != "ERROR" {
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
	glog.SetItemsList([]KeyId{KeyVersion})

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")
	if dl[0] != "1.0.0" {
		t.Error("Error version in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestHostname(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	glog.SetItemsList([]KeyId{KeyHostName})
	hostname := "hogeserver"
	glog.Config.st.hostname = hostname

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")
	if dl[0] != hostname {
		t.Error("Error hostname in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestUsername(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	glog.SetItemsList([]KeyId{KeyUserName})
	username := "suganoo"
	glog.Config.st.username = username

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")
	if dl[0] != username {
		t.Error("Error username in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestMessage(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	glog.SetItemsList([]KeyId{KeyMessage})
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
			if dl[0] != msg {
				t.Error("Error string in log.")
			}
		case 1:
			if dl[0] != strconv.Itoa(num) {
				t.Error("Error int in log.")
			}
		case 2:
			if dl[0] != fmt.Sprintf("%v", hs) {
				t.Error("Error struct in log.")
			}
		case 3:
			if dl[0] != (msg + " " + strconv.Itoa(num)) {
				t.Error("Error multi in log.")
			}
		case 4:
			if dl[0] != fmt.Sprintf("%v", list) {
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
	glog.SetItemsList([]KeyId{KeyProcessId})

	r := regexp.MustCompile(`\d`)

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")

	if ! r.MatchString(dl[0]) {
		t.Error("Error process id in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestGoroutineId(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	glog.SetItemsList([]KeyId{KeyGoroutineId})

	r := regexp.MustCompile(`gid:\d`)

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	dl := strings.Split(string(data), "\t")

	if ! r.MatchString(dl[0]) {
		t.Error("Error goroutine id in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestFileName(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	glog.SetItemsList([]KeyId{KeyFileName})

	r := regexp.MustCompile(`[gologger_test.go:\d]`)

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))
	if ! r.MatchString(string(data)) {
		t.Error("Error filename in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}

func TestFunc(t *testing.T) {
	os.Remove(testlog)
	glog := NewGologger(Configuration{Logfile: testlog})
	glog.SetItemsList([]KeyId{KeyFunc})

	glog.Info()

	data, _ := ioutil.ReadFile(testlog)
	t.Log(string(data))

	if strings.Split(string(data), "\t")[0] != "TestFunc" {
		t.Error("Error function name in log.")
	}
	
	glog.CloseFile()
	os.Remove(testlog)
}
