package ftl

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
)

type Log interface {
	BeginStep(name string)
	EndStep(State)
	Error(v ...interface{})
	Info(v ...interface{})
	Debug(v ...interface{})
	Trace(v ...interface{})
}

type Ops struct {
	Log Log
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

type StepFunc func(*Ops) State

func (op *Ops) Step(name string, f StepFunc) {
	op.Log.BeginStep(name)
	state := f(op)
	op.Log.EndStep(state)
}

func (op *Ops) Error(opName string, err error) {
	op.Log.Error(opName, err)
	debug.PrintStack()
	os.Exit(1)
}

func (op *Ops) Installed(name string) bool {
	cmd := exec.Command("dpkg-query", "--show", "--showformat='${db:Status-Status}'", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			if exit.ExitCode() == 1 {
				return false
			}
		}

		op.Error("Ops.Installed", err)
		return false
	}

	return string(out) == "'installed'"
}

func (op *Ops) UpdateRepos() {
	cmd := exec.Command("apt-get", "update")
	if err := cmd.Run(); err != nil {
		op.Error("Ops.UpdateRepos", err)
		return
	}
}

func (op *Ops) MissingPackage(name string) bool {
	return len(op.MissingPackages(name)) > 0
}

func (op *Ops) MissingPackages(name ...string) []string {
	var missing []string
	for _, n := range name {
		if !op.Installed(n) {
			missing = append(missing, n)
		}
	}
	return missing
}

func (op *Ops) Install(names ...string) {
	op.Log.Debug("Ops.Install:", strings.Join(names, ", "))
	cmd := exec.Command("apt-get", append([]string{"install", "--yes"}, names...)...)
	if err := cmd.Run(); err != nil {
		op.Error("Ops.Install", err)
		return
	}
}

func (op *Ops) DistroCodename() string {
	cmd := exec.Command("lsb_release", "--short", "--codename")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		op.Error("Ops.DistroCodename", err)
		return ""
	}
	return strings.Trim(out.String(), "\n")
}

func (op *Ops) AddRepo(repo, pubKey string) {
	resp, err := http.Get(pubKey)
	if err != nil {
		op.Error("Ops.AddRepo", err)
		return
	}
	defer resp.Body.Close()

	cmd := exec.Command("apt-key", "add", "-")
	cmd.Stdin = resp.Body
	if err := cmd.Run(); err != nil {
		op.Error("Ops.AddRepo", err)
		return
	}

	cmd = exec.Command("add-apt-repository", "--update", repo)
	if err := cmd.Run(); err != nil {
		op.Error("Ops.AddRepo", err)
		return
	}
}

func (op *Ops) ListRepos() []string {
	cmd := exec.Command("grep", "--recursive", "--no-filename", "--include", "*.list", "^deb ", "/etc/apt/sources.list", "/etc/apt/sources.list.d/")

	out, err := cmd.CombinedOutput()
	if err != nil {
		op.Error("Ops.AddRepo", err)
		return []string{}
	}

	var repos []string
	s := bufio.NewScanner(bytes.NewReader(out))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		repos = append(repos, s.Text())
	}
	return repos
}

func (op *Ops) MissingRepo(name string) bool {
	return len(op.MissingRepos(name)) > 0
}

func (op *Ops) MissingRepos(name ...string) []string {
	var missing []string
	existing := op.ListRepos()

	for _, n := range name {
		op.Log.Trace("Ops.MissingRepos: looking for:", n)
		var found bool
		for _, e := range existing {
			op.Log.Trace("Ops.MissingRepos: comparing with:", e)
			if n == e {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, n)
		}
	}

	return missing
}
