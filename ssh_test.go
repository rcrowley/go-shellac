package shellac

import (
	"github.com/rcrowley/go-shellac/ssh"
	"testing"
)

func TestSSH(t *testing.T) {
	testArgs(t, []string{}, Args(ssh.SSH{}))
	testArgs(t, []string{"example.com"}, Args(ssh.SSH{
		Hostname: "example.com",
	}))
}

func TestSSHOptions(t *testing.T) {
	testArgs(t, []string{"example.com"}, Args(ssh.SSH{
		Hostname: "example.com",
	}))
}
