# ftl
Configuration management in plain Go code.

**This is an early draft. Proceed with caution.**

With FTL you write Go programs that use a combination of `ftl`, `ftl/ops`, `ftl/log` and third-party libraries to update system state.

The `ftl` CLI copies and runs the plan on the destination host.

## Example

```
$ cat examples/tools/main.go
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
$ ftl root@172.13.0.1 examples/tools
2021/04/06 00:01:29 ((( [install tools]
2021/04/06 00:01:29 DBG [install tools] ops.UpdateRepos
2021/04/06 00:01:32 DBG [install tools] ops.Install: tree
2021/04/06 00:01:35 ))) [install tools] -> changed
$ ftl root@172.13.0.1 examples/tools
2021/04/06 00:02:36 ((( [install tools]
2021/04/06 00:02:36 ))) [install tools] -> unchanged
```
