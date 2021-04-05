package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	syncLocalLib = flag.String("sync-local-lib", "./", "copy local library to remote host")
)

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("usage: ftl [OPTIONS] user@target PLAN")
		fmt.Println("\nOPTIONS are:")
		flag.PrintDefaults()
		fmt.Println("\nPLAN is:")
		fmt.Println("  A path to a go file: copies the file to the remote host and runs it with `go run`")
		fmt.Println("  A directory: copies the directory to the remote host and runs `go run main.go` inside")
		os.Exit(1)
	}

	dest := flag.Arg(0)
	planPath := flag.Arg(1)

	if *syncLocalLib != "" {
		findRemoteHome := exec.Command("ssh", dest, "echo $HOME")
		remoteHome := strings.Trim(mustSucceed(findRemoteHome), "\n")

		target := fmt.Sprintf("%s/go/src/github.com/ftlops/ftl/", remoteHome)
		createTarget := exec.Command("ssh", dest, fmt.Sprintf("mkdir -p %s", target))
		mustSucceed(createTarget)

		remoteTarget := fmt.Sprintf("%s:%s", dest, target)

		// sync respecting .gitignore and removing locally deleted files from remote host
		syncLib := exec.Command("rsync", "-vha", *syncLocalLib, remoteTarget, "--include=**.gitignore", "--exclude=/.git", "--filter=:- .gitignore", "--delete-after")
		log.Printf("syncing library\n%s", mustSucceed(syncLib))
	}

	fi, err := os.Stat(planPath)
	if err != nil {
		log.Fatal(err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	remotePlanPath := fmt.Sprintf("/var/local/ftlops/ftl/plans/%s/", timestamp)

	mustSucceed(exec.Command("ssh", dest, fmt.Sprintf("mkdir -p %s", remotePlanPath)))

	if fi.Mode().IsDir() {
		// rsync needs the trailing slash to copy dir content and not dir itself
		if planPath[len(planPath)-1] != '/' {
			planPath += "/"
		}
	}
	mustSucceed(exec.Command("rsync", "-vha", planPath, fmt.Sprintf("%s:%s", dest, remotePlanPath)))

	var runCmd *exec.Cmd
	if fi.Mode().IsDir() {
		runCmd = exec.Command("ssh", dest, fmt.Sprintf("cd %s && go run %smain.go", remotePlanPath, remotePlanPath))
	} else {
		runCmd = exec.Command("ssh", dest, fmt.Sprintf("go run %s", remotePlanPath))
	}

	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	if err := runCmd.Run(); err != nil {
		log.Println("remote plan failed:", err)
	}
}

func mustSucceed(cmd *exec.Cmd) string {
	log.Println(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Oops, something went wrong.\n\n\tCommand: %s\n\tError: %s\n\tOutput:\n\n%s\n", cmd.String(), err, string(out))
		os.Exit(1)
	}
	return string(out)
}
