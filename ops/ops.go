package ops

import (
	"bufio"
	"bytes"
	"net/http"
	"os/exec"
	"strings"
)

type Log interface {
	Fatal(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
	Trace(v ...interface{})
}

var Logger Log

func Installed(name string) bool {
	cmd := exec.Command("dpkg-query", "--show", "--showformat='${db:Status-Status}'", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			if exit.ExitCode() == 1 {
				return false
			}
		}

		Logger.Fatal("ops.Installed:", err)
	}

	return string(out) == "'installed'"
}

func UpdateRepos() {
	Logger.Debug("ops.UpdateRepos")
	cmd := exec.Command("apt-get", "update")
	if err := cmd.Run(); err != nil {
		Logger.Fatal("ops.UpdateRepos:", err)
	}
}

func MissingPackage(name string) bool {
	return len(MissingPackages(name)) > 0
}

func MissingPackages(name ...string) []string {
	var missing []string
	for _, n := range name {
		if !Installed(n) {
			missing = append(missing, n)
		}
	}
	return missing
}

func Install(names ...string) {
	Logger.Debug("ops.Install:", strings.Join(names, ", "))
	cmd := exec.Command("apt-get", append([]string{"install", "--yes"}, names...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		Logger.Debug("ops.Install: combined output:\n\n", string(out))

		Logger.Fatal("ops.Install:", err)
	}
}

func DistroCodename() string {
	cmd := exec.Command("lsb_release", "--short", "--codename")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		Logger.Fatal("ops.DistroCodename:", err)
	}
	return strings.Trim(out.String(), "\n")
}

func AddRepo(repo, pubKey string) {
	resp, err := http.Get(pubKey)
	if err != nil {
		Logger.Fatal("ops.AddRepo:", err)
	}
	defer resp.Body.Close()

	cmd := exec.Command("apt-key", "add", "-")
	cmd.Stdin = resp.Body
	if err := cmd.Run(); err != nil {
		Logger.Fatal("ops.AddRepo:", err)
	}

	cmd = exec.Command("add-apt-repository", "--update", repo)
	if err := cmd.Run(); err != nil {
		Logger.Fatal("ops.AddRepo:", err)
	}
}

func ListRepos() []string {
	cmd := exec.Command("grep", "--recursive", "--no-filename", "--include", "*.list", "^deb ", "/etc/apt/sources.list", "/etc/apt/sources.list.d/")

	out, err := cmd.CombinedOutput()
	if err != nil {
		Logger.Fatal("ops.ListRepos:", err)
	}

	var repos []string
	s := bufio.NewScanner(bytes.NewReader(out))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		repos = append(repos, s.Text())
	}
	return repos
}

func MissingRepo(name string) bool {
	return len(MissingRepos(name)) > 0
}

func MissingRepos(name ...string) []string {
	var missing []string
	existing := ListRepos()

	for _, n := range name {
		Logger.Trace("ops.MissingRepos: looking for:", n)
		var found bool
		for _, e := range existing {
			Logger.Trace("ops.MissingRepos: comparing with:", e)
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
