package ftl

var DefaultOps = Ops{&LogLog{}}

func Step(name string, f StepFunc) {
	DefaultOps.Step(name, f)
}
