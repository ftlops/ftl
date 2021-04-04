package ftl

import (
	"fmt"

	"github.com/ftlops/ftl/log"
	"github.com/ftlops/ftl/ops"
)

func init() {
	ops.Logger = &log.DefaultLogger
}

type State int

const (
	StateUnchanged State = iota
	StateChanged
)

func (s State) String() string {
	switch s {
	case StateUnchanged:
		return "unchanged"
	case StateChanged:
		return "changed"
	default:
		return fmt.Sprintf("unknown (%d)", s)
	}
}

type StepFunc func() State

func Step(name string, f StepFunc) {
	log.BeginStep(name)
	state := f()
	log.EndStep(state.String())
}
