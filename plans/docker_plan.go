package main

import (
	"fmt"
	"log"

	"github.com/ftlops/ftl"
)

func main() {
	installDocker()
}

func installDocker() {
	log.Println(ftl.ListRepos())

	log.Println("* install docker")
	if ftl.MissingPackage("docker-ce") {
		prereqs := []string{"apt-transport-https", "ca-certificates", "gnupg", "lsb-release"}
		if missing := ftl.MissingPackages(prereqs...); len(missing) > 0 {
			log.Println("** missing: ", missing)
			log.Println("** update repos")
			ftl.UpdateRepos()
			log.Println("** install prerequisites")
			ftl.Install(missing...)
		}

		log.Println("** get distro codename")
		codename := ftl.DistroCodename()

		repo := fmt.Sprintf("deb [arch=amd64] https://download.docker.com/linux/ubuntu %s stable", codename)
		if ftl.MissingRepo(repo) {
			log.Println("** add repo")
			ftl.AddRepo(repo, "https://download.docker.com/linux/ubuntu/gpg")
		}
		log.Println("** install docker-ce")
		ftl.Install("docker-ce")
	}
}
