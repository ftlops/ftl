package main

import (
	"fmt"
	"log"

	"github.com/ngrash/ftl"
)

func main() {
	installDocker()
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func installDocker() {
	if len(ftl.Missing("docker-ce")) > 0 {
		prereqs := []string{"apt-transport-https", "ca-certificates", "curl", "gnupg", "lsb-release"}
		if missing := ftl.Missing(prereqs...); len(missing) > 0 {
			must(ftl.UpdateRepos())
			must(ftl.Install(missing...))
		}

		codename, err := ftl.DistroCodename()
		if err != nil {
			log.Fatal(err)
		}
		must(ftl.AddRepo(fmt.Sprintf("deb [arch=amd64] https://download.docker.com/linux/ubuntu %s stable", codename), "https://download.docker.com/linux/ubuntu/gpg"))
		must(ftl.UpdateRepos())
		must(ftl.Install("docker-ce"))
	}
}
