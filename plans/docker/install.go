package docker

import (
	"fmt"

	"github.com/ftlops/ftl"
	"github.com/ftlops/ftl/ops"
)

func Install() {
	ftl.Step("install docker", func() ftl.State {
		if !ops.MissingPackage("docker-ce") {
			return ftl.StateUnchanged
		}

		ftl.Step("install prereqs", func() ftl.State {
			prereqs := []string{
				"apt-transport-https",
				"ca-certificates",
				"gnupg",
				"lsb-release",
			}
			missing := ops.MissingPackages(prereqs...)
			if len(missing) == 0 {
				return ftl.StateUnchanged
			}

			ops.UpdateRepos()
			ops.Install(missing...)
			return ftl.StateChanged
		})

		ftl.Step("add repo", func() ftl.State {
			codename := ops.DistroCodename()
			repo := fmt.Sprintf("deb [arch=amd64] https://download.docker.com/linux/ubuntu %s stable", codename)
			if !ops.MissingRepo(repo) {
				return ftl.StateUnchanged
			}
			ops.AddRepo(repo, "https://download.docker.com/linux/ubuntu/gpg")
			return ftl.StateChanged
		})

		ops.Install("docker-ce")
		return ftl.StateChanged
	})
}
