package shellac

import (
	"fmt"
	"github.com/rcrowley/go-shellac/coreutils"
	"github.com/rcrowley/go-shellac/ssh"
	"os/user"
)

func ExampleFind() {
	Run(coreutils.Find{
		Dirnames: []string{"."},
		Name:     "*.go",
		Type:     coreutils.FindFile,
	})
}

func ExampleFindChannel() {
	ch := make(chan string)
	go func() {
		cmd := Command(coreutils.Find{
			Dirnames: []string{"."},
			Name:     "*.go",
			Type:     coreutils.FindFile,
		})
		cmd.ChannelStdout(ch)
		cmd.Run()
	}()
	for s := range ch {
		fmt.Println(s)
	}
}

func ExampleSSH() {
	u, _ := user.Current()
	Run(ssh.SSH{
		AgentForwarding: true,
		Command:         []string{"hostname"},
		Hostname:        "example.com",
		Login:           u.Username,
		Options:         ssh.SSHOptions{"StrictHostKeyChecking": "yes"},
	})
}

func ExampleSudoFind() {
	Sudo(coreutils.Find{
		Dirnames: []string{"/dev"},
		Readable: true,
		Type:     coreutils.FindBlock,
	})
}
