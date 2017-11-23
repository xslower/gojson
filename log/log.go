package log

import (
	"os"
	"path/filepath"
	"strings"
	"io"
	"github.com/resure-tech/lib/base/nolock"
	"time"
	"github.com/resure-tech/lib/goutils/utils"
	"log"
	"github.com/resure-tech/lib/base/types"
)

const (
	//日志级别，越高信息越少
	DEBUG   uint8 = iota
	INFO
	WARNING
	FATAL
)

var (
	_mode_map = map[string]uint8{
		`debug`:   DEBUG,
		`info`:    INFO,
		`warning`: WARNING,
		`fatal`:   FATAL,
		`0`:       DEBUG,
		`1`:       INFO,
		`2`:       WARNING,
		`3`:       FATAL,
	}
)

func StrToMode(mode string) uint8 {
	m := strings.ToLower(mode)
	i := _mode_map[m]
	return i
}

type Config struct {
	FileName    string
	LogMode     string
	FileMaxLine uint32
}

func NewLogger(conf Config) (l *Logger) {
	if conf.FileName == `` {
		panic(`log file path is empty`)
	}
	var dir = filepath.Dir(conf.FileName)
	_ = os.MkdirAll(dir, 0777)
	l = &Logger{
		mode:    StrToMode(conf.LogMode),
		file:    conf.FileName,
		lock:    nolock.NewLocker(),
		maxLine: conf.FileMaxLine}
	if conf.FileMaxLine > 1000 {
		l.NewLogFile(true)
		l.countFunc = l.Count
	} else {
		l.NewLogFile(false)
		l.countFunc = func() {}
	}
	return
}

/**
go的日志输出性能主要消耗在使用反正把对象序列化上。这里的包只接收string。
//todo 自己写文件，代替logger包。
 */
type Logger struct {
	log  *log.Logger
	file string
	lock *nolock.Locker
	w    io.WriteCloser
	//f       *os.File
	counter   uint32
	maxLine   uint32
	countFunc func()
	mode      uint8
}

func (this *Logger) SetMode(mode uint8) {
	this.mode = mode
}

func (this *Logger) Debug(msg ...string) {
	if this.mode > DEBUG {
		return
	}
	this.PrefixPrintln(`[Debug]`, msg...)
}

func (this *Logger) Info(msg ...string) {
	if this.mode > INFO {
		return
	}
	this.PrefixPrintln(`[Info]`, msg...)
}

func (this *Logger) Warning(msg ...string) {
	if this.mode > WARNING {
		return
	}
	this.PrefixPrintln(`[Warning]`, msg...)
}

func (this *Logger) Fatal(msg ...string) {
	this.PrefixPrintln(`[Fatal]`, msg...)
}

func (this *Logger) PrintBytes(bmsg []byte) {
	defer this.lock.WUnlock()
	this.lock.WLock()
	this.w.Write(bmsg)
	this.countFunc()
}

func (this *Logger) PrefixPrintln(prefix string, msg ...string) {
	buf := types.NewBufBytes(len(msg) * 30)
	now := utils.StdDateTime(time.Now())
	buf.WriteString(now)
	buf.WriteByte(' ')
	buf.WriteString(prefix)
	for i := 0; i < len(msg); i++ {
		buf.WriteByte(' ')
		buf.WriteString(msg[i])
	}
	buf.WriteByte('\n')

	this.PrintBytes(buf.Bytes())
}

func (this *Logger) Count() {
	this.counter++
	if this.counter > this.maxLine {
		this.NewLogFile(true)
		this.counter = 0
	}
}

func (this *Logger) NewLogFile(newPart bool) {
	defer this.lock.WUnlock()
	this.lock.WLock()
	if this.w != nil {
		this.w.Close()
	}
	fn := this.file
	if newPart {
		fn += `_` + utils.Date(`Y-m-d_H-i-s`, time.Now())
	}
	var err error
	this.w, err = os.OpenFile(fn, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err.Error())
	}
	return
}
