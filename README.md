Shellac
=======

A declarative, strongly-typed API for executing shell commands from Go.

Installation
------------

```sh
go get "github.com/rcrowley/go-shellac"
```

Usage
-----

Define `struct`s tagged so `shellac.Args` can construct commands:

```go
type Find struct {
	Dirnames []string `pos:"first"`
	Mode *int `flag:"-perm" format:"%o"`
	Name string `flag:"-name"`
    // ...
}
```

See <http://godoc.org/github.com/rcrowley/go-shellac#Args> for complete details on defining these `struct`s.

Then use your `struct` or one distributed with Shellac to execute shell commands safely from Go:

```go
import (
    "github.com/rcrowley/go-shellac"
    "github.com/rcrowley/go-shellac/coreutils"
    "github.com/rcrowley/go-shellac/ssh"
)

shellac.Run(coreutils.Find{
    Dirnames: []string{"."},
    Name: "*.go",
    Type: coreutils.FindFile,
})

shellac.Run(ssh.SSH{
    AgentForwarding: true,
    Command: []string{"hostname"},
    Hostname: "example.com",
    Login: "example",
    Options: ssh.SSHOptions{"StrictHostKeyChecking": "yes"},
})
```

TODO
----

* Support for the rest of GNU coreutils, `scp`, and whatever else we need.
* Create a tool that turns a `flag.FlagSet` into a Shellac-compatible `struct`.

TODONE
------

* `struct`-to-`exec.Cmd` parser.
* Field ordering via tags.
* `sudo`(8) support.
* Custom string formatters.
* Support for multiple arguments per option.
* `find`(1) implementation in the `coreutils` package.
* `ssh`(1) implementation in the `ssh` package.
* Package documentation.
* Channels for standard input, output, and error.
