package main

import (
	"fmt"

	"github.com/ftlops/ftl"
)

func main() {
	installDocker()
}

func installDocker() {

	ftl.Step("install docker", func(op *ftl.Ops) ftl.State {
		if !op.MissingPackage("docker-ce") {
			return ftl.StateUnchanged
		}

		op.Step("install prereqs", func(op *ftl.Ops) ftl.State {
			prereqs := []string{
				"apt-transport-https",
				"ca-certificates",
				"gnupg",
				"lsb-release",
			}
			missing := op.MissingPackages(prereqs...)
			if len(missing) == 0 {
				return ftl.StateUnchanged
			}

			op.Log.Info("update repos")
			op.UpdateRepos()
			op.Log.Info("install prerequisites")
			op.Install(missing...)
			return ftl.StateChanged
		})

		op.Step("add repo", func(op *ftl.Ops) ftl.State {
			codename := op.DistroCodename()
			repo := fmt.Sprintf("deb [arch=amd64] https://download.docker.com/linux/ubuntu %s stable", codename)
			if !op.MissingRepo(repo) {
				return ftl.StateUnchanged
			}
			op.AddRepo(repo, "https://download.docker.com/linux/ubuntu/gpg")
			return ftl.StateChanged
		})

		op.Log.Info("install docker-ce")
		op.Install("docker-ce")
		return ftl.StateChanged
	})
}
