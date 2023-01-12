package ulog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

func init() {

	// 监听系统关闭之前
	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {

			// 等待call信号 执行逻辑
			<-closeCh
			// "停止运行了,可以执行释放资源的方法 "
		}
	}()
}

// LoggerBuilder 日志生成器
type LoggerBuilder struct {
	FilePath string    // FilePath 文件路径
	Out      io.Writer // Out 控制台的Writer对象,可以通过 os.Stdout 获取
}

// loggerFile 日志文件
type loggerFile struct {
	fullPath  string // fullPath 全路径
	fileName  string // fileName 文件名字
	baiscPath string // baiscPath 基础路径
}

// spInfo 拆分路径信息
func (logFile *loggerFile) spPathInfo() {
	fmt.Println("fullPath: ", logFile.fullPath)
	idx := strings.LastIndex(logFile.fullPath, "/")
	r := []rune(logFile.fullPath)
	logFile.baiscPath = string(r[:idx])
	logFile.fileName = string(r[idx+1:])
}

type Logger struct {
	// 读写
	writer []io.Writer

	// 读取暂时不再考虑范围内
	// reader *io.Reader
}

var logger *Logger

// New 创建一个新的日志
func New(builder *LoggerBuilder) *Logger {
	writers := make([]io.Writer, 0)

	// 如果没有传文件路径,就不写文件
	if len(builder.FilePath) > 0 {
		info := &loggerFile{fullPath: builder.FilePath}
		info.spPathInfo()
		fmt.Printf("%v\n", info)
		file := fileCreate(info)
		writers = append(writers, file)
	}

	if builder.Out != nil {
		writers = append(writers, builder.Out)
	}

	logger = &Logger{writer: writers}
	return logger
}

const (
	InfoLevel  = "info"
	DebugLevel = "debug"
	ErrorLevel = "error"
)

func (log *Logger) Info(msg string) {
	c := callerInfoBuilder(runtime.Caller(0))
	c.level = InfoLevel
	write(log.writer, msg, c)
}

func (log *Logger) Infof(f, msg string) {
	log.Info(fmt.Sprintf(f, msg))
}

func (log *Logger) Debug(msg string) {
	c := callerInfoBuilder(runtime.Caller(0))
	c.level = DebugLevel
	write(log.writer, msg, c)
}

func (log *Logger) Debugf(f, msg string) {
	log.Debug(fmt.Sprintf(f, msg))
}

func (log *Logger) Error(msg string) {
	c := callerInfoBuilder(runtime.Caller(0))
	c.level = ErrorLevel
	write(log.writer, msg, c)
}

func (log *Logger) Errorf(f, msg string) {
	log.Error(fmt.Sprintf(f, msg))
}

// write 写文件
func write(ws []io.Writer, msg string, cinfo callerInfo) {
	executeEvent(getEventMonitor(BeforeEvent), msg)

	wg := sync.WaitGroup{}
	for _, w := range ws {
		wg.Add(1)
		go func(writer io.Writer) {

			// 头部信息
			head, _ := colorFiexdhead(cinfo)

			buffer := bufio.NewWriter(writer)
			buffer.WriteString(fmt.Sprintf("%v %v\n", head, msg))
			buffer.Flush()
			wg.Done()
		}(w)

		executeEvent(getEventMonitor(ProcessEvent), msg)
	}
	wg.Wait()

	executeEvent(getEventMonitor(AfterEvent), msg)
}

// 继承自io.Writer
func (log *Logger) Writer(b []byte) (n int, err error) {
	if len(b) <= -1 {
		return 0, errors.New("bytes 长度必须大于0")
	}
	c := callerInfoBuilder(runtime.Caller(1))
	c.level = DebugLevel
	write(log.writer, string(b), c)
	return len(b), nil
}

// fileCreate 创建文件
func fileCreate(loggerFile *loggerFile) (file *os.File) {
	defer file.Close()

	// 支持日志分隔

	fullPath := loggerFile.fullPath
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(loggerFile.baiscPath, os.ModePerm)
		// 创建这个文件
		file, err = os.Create(loggerFile.fullPath)
		if err != nil {
			log.Panic(err)
		}
	} else {

		// 文件存在的时候直接打开
		file, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}

	}

	return file
}

func executeEvent(events []*event, msg string) {
	for _, e := range events {
		e.execute(msg)
	}
}
