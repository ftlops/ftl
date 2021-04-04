package ftl

import (
	"bytes"
	"net/http"
	"os/exec"
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

func Installed(name string) (bool, error) {
	_, err := Exec("dpkg", "-l", name)
	if exit, ok := err.(*exec.ExitError); ok {
		if exit.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func UpdateRepos() error {
	_, err := Exec("apt-get", "update")
	return err
}

func Missing(names ...string) []string {
	var missing []string
	for _, name := range names {
		installed, err := Installed(name)
		if err != nil {
			panic(err)
		}
		if !installed {
			missing = append(missing, name)
		}
	}
	return missing
}

func Install(names ...string) error {
	_, err := Exec("apt-get", append([]string{"install", "--yes"}, names...)...)
	if err != nil {
		return err
	}

	return nil
}

func DistroCodename() (string, error) {
	cmd, err := Exec("lsb_release", "--short", "--codename")
	if err != nil {
		return "", err
	}
	return cmd.Stdout.(*bytes.Buffer).String(), nil
}

func AddRepo(repo, pubKey string) error {
	resp, err := http.Get(pubKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cmd := exec.Command("apt-key", "add", "-")
	cmd.Stdin = resp.Body
	if err := cmd.Run(); err != nil {
		return err
	}

	_, err = Exec("add-apt-repository", repo)
	if err != nil {
		return err
	}

	return nil
}
