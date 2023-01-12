package ulog

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	defaultPath = "logs/tp1.log"
)

func TestLogger(t *testing.T) {
	log := New(&LoggerBuilder{FilePath: defaultPath, Out: os.Stdout})
	log.Info("test")
}

// TestWriteFile 测试写文件的情况
func TestWriteFile(t *testing.T) {
	// 真实的file文件
	log := New(&LoggerBuilder{FilePath: defaultPath})

	// 设置事件
	unix := time.Now().Unix()
	EventBuilder(BeforeEvent, func(et *EventType, s string) {
		fmt.Println(unix, " - ", s)
	})

	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(l *Logger, i int) {
			l.Info(fmt.Sprintln("这是一条测试日志 😄 i=", i))
			wg.Done()
		}(log, i)
	}
	wg.Wait()

}

// TestLoggerFormat 日志格式化测试
func TestLoggerFormat(t *testing.T) {
	log := New(&LoggerBuilder{Out: os.Stdout})
	log.Info("测试日志信息")
	log.Debug("警告信息")
	log.Error("错误信息")
}
