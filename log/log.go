package log

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

var DefaultLogger = Logger{}

func BeginStep(name string)  { DefaultLogger.BeginStep(name) }
func EndStep(result string)  { DefaultLogger.EndStep(result) }
func Fatal(v ...interface{}) { DefaultLogger.Fatal(v...) }
func Error(v ...interface{}) { DefaultLogger.Error(v...) }
func Info(v ...interface{})  { DefaultLogger.Info(v...) }
func Debug(v ...interface{}) { DefaultLogger.Debug(v...) }
func Trace(v ...interface{}) { DefaultLogger.Trace(v...) }

type Logger struct {
	stack   []string
	lastFul string
}

func (l *Logger) BeginStep(name string) {
	l.stack = append(l.stack, name)

	l.printf("(((", []interface{}{})
}

func (l *Logger) EndStep(result string) {
	l.printf(")))", []interface{}{"->", result})
	l.stack = l.stack[:len(l.stack)-1]
}

func (l *Logger) Fatal(v ...interface{}) {
	l.printf("FTL", v)
	debug.PrintStack()
	os.Exit(1)
}

func (l *Logger) Error(v ...interface{}) {
	l.printf("ERR", v)
}

func (l *Logger) Info(v ...interface{}) {
	l.printf("INF", v)
}

func (l *Logger) Debug(v ...interface{}) {
	l.printf("DBG", v)
}

func (l *Logger) Trace(v ...interface{}) {
	l.printf("TRC", v)
}

func (l *Logger) printf(lvl string, v []interface{}) {
	full := l.fmtStack()
	if l.lastFul != full {
		log.Println(l.fmt(lvl, full, v)...)
		l.lastFul = full
	} else {
		shortStack := "["
		prefix := "... > "
		basename := l.stack[len(l.stack)-1]
		padLeft := len(full) - len(basename) - len(prefix) - 2
		for i := 0; i < padLeft; i++ {
			shortStack += " "
		}
		if padLeft > 0 {
			shortStack += prefix
		}
		shortStack += basename + "]"
		log.Println(l.fmt(lvl, shortStack, v)...)
	}
}

func (l *Logger) fmt(lvl, stack string, v []interface{}) []interface{} {
	var vs []interface{}
	vs = append(vs, lvl)
	vs = append(vs, stack)
	for _, x := range v {
		vs = append(vs, x)
	}
	return vs
}

func (l *Logger) fmtStack() string {
	return fmt.Sprintf("[%s]", strings.Join(l.stack, " > "))
}
