package config

import (
	"Tally/global"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

type Log struct{}

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

	time := f.Time.Format("2006-01-02 15:04:05")
	if f.HasCaller() {
		funcpVal := f.Caller.Function
		fileval := fmt.Sprintf("%s:%d", path.Base(f.Caller.Function), f.Caller.Line)
		fmt.Fprintf(b, "[%s] [%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", Config.Logs.Prefix, time, leave, f.Level, fileval, funcpVal, f.Message)
	} else {
		fmt.Fprintf(b, "[%s] [%s] \x1b[%dm[%s]\x1b[0m %s\n", Config.Logs.Prefix, time, leave, f.Level, f.Message)
	}
	return b.Bytes(), nil

}

func InitLog() {
	m := log.New()

	//自定义输出
	m.SetFormatter(&Log{})

	//输出到控制台
	m.SetOutput(os.Stdout)
	//输出任务和行号
	m.SetReportCaller(true)
	//最低输出级别
	m.SetLevel(log.InfoLevel)
	global.Global.Log = m
}
