package config

import (
	"Tally/global"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

type Log struct {
	w io.Writer
	m int
}

func (l *Log) Format(f *log.Entry) ([]byte, error) {
	var leave int
	switch f.Level {
	case log.InfoLevel, log.DebugLevel:
		leave = global.Gray
	case log.WarnLevel:
		leave = global.Yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		leave = global.Red
	default:
		leave = global.Blue
	}
	var b *bytes.Buffer
	if f.Buffer != nil {
		b = f.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	time1 := f.Time.Format("2006-01-02 15:04:05")
	if f.HasCaller() {
		funcpVal := f.Caller.Function
		fileval := fmt.Sprintf("%s:%d", path.Base(f.Caller.Function), f.Caller.Line)
		fmt.Fprintf(b, "[%s] [%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", Config.Logs.Prefix, time1, leave, f.Level, fileval, funcpVal, f.Message)
	} else {
		fmt.Fprintf(b, "[%s] [%s] \x1b[%dm[%s]\x1b[0m %s\n", Config.Logs.Prefix, time1, leave, f.Level, f.Message)
	}
	return b.Bytes(), nil

}

// 这个是重写
func (l *Log) Write(p []byte) (n int, err error) {
	n, err = l.w.Write(p)
	l.m += n
	return n, err
}

func (l *Log) logFile(m *log.Logger) {
	t := time.Now().Format(time.DateTime)
	//创建文件
	file, err := os.Create(t + ".log")
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	s := time.NewTicker(time.Minute * 60)
	ls := new(Log)
	ls.w = file
	//输出到控制台
	m.SetOutput(io.MultiWriter(os.Stdout, ls.w))
	for {
		select {
		case <-s.C:
			//	判读是否超过100m
			if l.m > 100*(1024*1024) {
				t = time.Now().Format(time.DateTime)
				file, err = os.Create(t + ".log")
				ls = new(Log)
				ls.w = file
				//输出到控制台,日志文件中
				m.SetOutput(io.MultiWriter(os.Stdout, ls.w))
			}
		}
	}

}

func InitLog() {
	m := log.New()

	//自定义输出
	m.SetFormatter(&Log{})

	//输出任务和行号
	m.SetReportCaller(true)
	//最低输出级别
	m.SetLevel(log.InfoLevel)
	global.Global.Log = m
}
