package main

import (
	"github.com/ftlops/ftl"
	"github.com/ftlops/ftl/ops"
)

func main() {
	ftl.Step("install tools", func() ftl.State {
		missing := ops.MissingPackages("gnupg", "tree", "htop")
		if len(missing) == 0 {
			return ftl.StateUnchanged
		}
		ops.UpdateRepos()
		ops.Install(missing...)
		return ftl.StateChanged
	})
}
