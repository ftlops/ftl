package ftl

import (
	"fmt"
	"log"
	"strings"
)

type LogLog struct {
	stack         []string
	lastFullStack string
}

func (ll *LogLog) BeginStep(name string) {
	ll.stack = append(ll.stack, name)

	ll.printf("(((", []interface{}{})
}

func (ll *LogLog) EndStep(st State) {
	ll.printf(")))", []interface{}{"->", st})
	ll.stack = ll.stack[:len(ll.stack)-1]
}

func (ll *LogLog) Error(v ...interface{}) {
	ll.printf("ERR", v)
}

func (ll *LogLog) Info(v ...interface{}) {
	ll.printf("INF", v)
}

func (ll *LogLog) Debug(v ...interface{}) {
	ll.printf("DBG", v)
}

func (ll *LogLog) Trace(v ...interface{}) {
	ll.printf("TRC", v)
}

func (ll *LogLog) printf(lvl string, v []interface{}) {
	fullStack := ll.fmtStack()
	if ll.lastFullStack != fullStack {
		log.Println(ll.fmt(lvl, fullStack, v)...)
		ll.lastFullStack = fullStack
	} else {
		shortStack := "["
		prefix := "... > "
		basename := ll.stack[len(ll.stack)-1]
		padLeft := len(fullStack) - len(basename) - len(prefix) - 2
		for i := 0; i < padLeft; i++ {
			shortStack += " "
		}
		if padLeft > 0 {
			shortStack += prefix
		}
		shortStack += basename + "]"
		log.Println(ll.fmt(lvl, shortStack, v)...)
	}
}

func (ll *LogLog) fmt(lvl, stack string, v []interface{}) []interface{} {
	var vs []interface{}
	vs = append(vs, lvl)
	vs = append(vs, stack)
	vs = append(vs, v...)
	return vs
}

func (ll *LogLog) fmtStack() string {
	return fmt.Sprintf("[%s]", strings.Join(ll.stack, " > "))
}
