package log

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	INFO int8 = iota
	DEBUG
	WARNING
	FATAL
)

var (
	_mode_map = map[string]int8{
		`info`:    INFO,
		`debug`:   DEBUG,
		`warning`: WARNING,
		`fatal`:   FATAL,
	}
)

func StrToMode(mode string) int8 {
	m := strings.ToLower(mode)
	i := _mode_map[m]
	return i
}

func NewLogger(filePath string, mode int8) (l *Logger) {
	var dir = filepath.Dir(filePath)
	_ = os.MkdirAll(dir, 0777)
	logFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("open file error !")
	}
	l = &Logger{mode: mode}
	l.log = log.New(logFile, "", log.LstdFlags)
	return
}

type Logger struct {
	mode int8
	log  *log.Logger
}

func (this *Logger) SetMode(mode int8) {
	this.mode = mode
}

func (this *Logger) Info(msg ...interface{}) {
	if this.mode > INFO {
		return
	}
	this.log.SetPrefix(`[Info] `)
	this.log.Println(msg...)
}

func (this *Logger) Debug(msg ...interface{}) {
	if this.mode > DEBUG {
		return
	}
	this.log.SetPrefix(`[Debug] `)
	this.log.Println(msg...)
}

func (this *Logger) Warning(msg ...interface{}) {
	if this.mode > WARNING {
		return
	}
	this.log.SetPrefix(`[Warning] `)
	this.log.Println(msg...)
}

func (this *Logger) Fatal(msg ...interface{}) {
	this.log.SetPrefix(`[Fatal] `)
	this.log.Println(msg...)
}
