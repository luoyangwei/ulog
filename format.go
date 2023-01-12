package ulog

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func init() {
	// 实际情况证明,获取时间的方式会耗时,所以给定了几个特定的时间格式
	// 在初始化的时候就把时间格式好,然后创建一个进程一直去修改时间内容
	dt = &DateTime{mutex: sync.RWMutex{}}
	go func() {

		for {
			currentTime := time.Now()
			dt.mutex.Lock()
			dt.ANSICVal = currentTime.Format(ANSIC)
			dt.UnixDateVal = currentTime.Format(UnixDate)
			dt.RubyDateVal = currentTime.Format(RubyDate)
			dt.RFC822Val = currentTime.Format(RFC822)
			dt.RFC822ZVal = currentTime.Format(RFC822Z)
			dt.RFC850Val = currentTime.Format(RFC850)
			dt.RFC1123Val = currentTime.Format(RFC1123)
			dt.RFC1123ZVal = currentTime.Format(RFC1123Z)
			dt.RFC3339Val = currentTime.Format(RFC3339)
			dt.CNRFC3339Val = currentTime.Format(CNRFC3339)
			dt.RFC3339NanoVal = currentTime.Format(RFC3339Nano)
			dt.CNRFC3339NanoVal = currentTime.Format(CNRFC3339Nano)
			dt.KitchenVal = currentTime.Format(Kitchen)
			dt.StampVal = currentTime.Format(Stamp)
			dt.StampMilliVal = currentTime.Format(StampMilli)
			dt.StampMicroVal = currentTime.Format(StampMicro)
			dt.StampNanoVal = currentTime.Format(StampNano)
			dt.mutex.Unlock()
		}

	}()

	go func() {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		sb := strings.Builder{}
		ss := strings.Split(dir, "/")
		for i, j := len(ss)-2, 0; i >= 0; i-- {
			sb.WriteString(fmt.Sprintf("/%s", ss[j]))
			j++
		}
		path := sb.String()
		basicWorkPath = path[1:]
		fmt.Println(basicWorkPath)
	}()
}

const (
	ANSIC         = "Mon Jan _2 15:04:05 2006"
	UnixDate      = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate      = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822        = "02 Jan 06 15:04 MST"
	RFC822Z       = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850        = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123       = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z      = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339       = "2006-01-02T15:04:05Z07:00"
	CNRFC3339     = "2006-01-02 15:04:05" // 国内经常会用到的格式
	RFC3339Nano   = "2006-01-02T15:04:05.999999999Z07:00"
	CNRFC3339Nano = "2006-01-02 15:04:05.999999999" // 国内经常会用到的格式
	Kitchen       = "3:04PM"
	Stamp         = "Jan _2 15:04:05"
	StampMilli    = "Jan _2 15:04:05.000"
	StampMicro    = "Jan _2 15:04:05.000000"
	StampNano     = "Jan _2 15:04:05.000000000"
)

var basicWorkPath string
var dt *DateTime

type DateTime struct {
	mutex sync.RWMutex

	ANSICVal         string
	UnixDateVal      string
	RubyDateVal      string
	RFC822Val        string
	RFC822ZVal       string
	RFC850Val        string
	RFC1123Val       string
	RFC1123ZVal      string
	RFC3339Val       string
	CNRFC3339Val     string
	RFC3339NanoVal   string
	CNRFC3339NanoVal string
	KitchenVal       string
	StampVal         string
	StampMilliVal    string
	StampMicroVal    string
	StampNanoVal     string
}

const (
	NotCallerInfoError = "无法获取调用者信息"
)

// callerInfo 调用者信息
type callerInfo struct {
	level string  // level 调用等级
	pc    uintptr // pc pc地址
	file  string  // file 调用路径
	line  int     // line 调用行
	ok    bool    // ok 是否获取成功
}

// callerInfoBuilder 生成调用者信息
func callerInfoBuilder(pc uintptr, file string, line int, ok bool) callerInfo {
	return callerInfo{
		pc:    pc,
		file:  file,
		line:  line,
		ok:    ok,
	}
}

func colorFiexdhead(caller callerInfo) (string, error) {
	if caller.ok {
		file := strings.Join([]string{strings.ReplaceAll(caller.file, basicWorkPath, ".."), fmt.Sprint(caller.line)}, ":")

		var sLevel string
		switch caller.level {
		case InfoLevel:
			sLevel = "\x1b[36m[%-6s]\x1b[0m"
		case DebugLevel:
			sLevel = "\x1b[33m[%-6s]\x1b[0m"
		case ErrorLevel:
			sLevel = "\x1b[31m[%-6s]\x1b[0m"
		default:
			// 默认跟info一样的
			sLevel = "\x1b[36m[%-6s]\x1b[0m"
		}

		sLevel = fmt.Sprintf(sLevel, strings.ToUpper(caller.level))
		s := fmt.Sprintf("\x1b[90m%s\x1b[0m [ %d ] - %s %-40s \x1b[36m->\x1b[0m ", dt.CNRFC3339Val, caller.pc, sLevel, file)
		return s, nil
	}
	return "", errors.New(NotCallerInfoError)
}

// type formatJson struct {
// 	Pc     uintptr `json:"pc"`
// 	File   string  `json:"file"`
// 	LineNo int     `json:"lineNo"`
// 	Msg    string  `json:"msg"`
// 	Date   string  `json:"date"`
// }

// func fJson(level string) (string, error) {
// 	pc, file, lineNo, ok := caller()
// 	if ok {
// 		sInfo := &formatJson{
// 			Pc:     pc,
// 			File:   file,
// 			LineNo: lineNo,
// 		}

// 		b, err := json.Marshal(sInfo)
// 		if err != nil {
// 			return "", err
// 		}
// 		return string(b), nil
// 	}
// 	return "", errors.New(NotCallerInfoError)
// }

// func colorFormat(fiexd, msg string) {

// }
