package ftl

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func Exec(command string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return cmd, err
	}
	return cmd, nil
}

func Installed(name string) bool {
	cmd := exec.Command("dpkg-query", "--show", "--showformat='${db:Status-Status}'", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			if exit.ExitCode() == 1 {
				return false
			}
		}
		panic(err)
	}

	return string(out) == "'installed'"
}

func UpdateRepos() {
	_, err := Exec("apt-get", "update")
	if err != nil {
		panic(err)
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
	_, err := Exec("apt-get", append([]string{"install", "--yes"}, names...)...)
	if err != nil {
		panic(err)
	}
}

func DistroCodename() string {
	cmd, err := Exec("lsb_release", "--short", "--codename")
	if err != nil {
		panic(err)
	}
	return strings.Trim(cmd.Stdout.(*bytes.Buffer).String(), "\n")
}

func AddRepo(repo, pubKey string) {
	resp, err := http.Get(pubKey)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cmd := exec.Command("apt-key", "add", "-")
	cmd.Stdin = resp.Body
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	_, err = Exec("add-apt-repository", "--update", repo)
	if err != nil {
		panic(err)
	}
}

func ListRepos() []string {
	cmd := exec.Command("grep", "--recursive", "--no-filename", "--include", "*.list", "^deb ", "/etc/apt/sources.list", "/etc/apt/sources.list.d/")

	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
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
		log.Printf("looking for: %s", n)
		var found bool
		for _, e := range existing {
			log.Printf("comparing with: %s", e)
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
