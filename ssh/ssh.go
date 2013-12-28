package ssh

import (
	"fmt"
	"sort"
	"strings"
)

// ssh(1)
type SSH struct {

	// -1
	SSHv1 bool `flag:"-1"`

	// -2
	SSHv2 bool `flag:"-2"`

	// -4
	IPv4 bool `flag:"-4"`

	// -6
	IPv6 bool `flag:"-6"`

	// -A
	AgentForwarding bool `flag:"-A"`

	// -a
	NoAgentForwarding bool `flag:"-a"`

	// -b <bind_address>
	BindAddress string `flag:"-b"`

	// -C
	Compression bool `flag:"-C"`

	// -c <cipher_spec>
	CipherSpec string `flag:"-c"`

	// -D [<bind_address>:]<port>
	DynamicForward string `flag:"-D"` // FIXME data structure + multiple

	// -e <escape_char>
	EscapeChar string `flag:"-e"`

	// -F <configfile>
	ConfigFile string `flag:"-F"`

	// -f
	Background bool `flag:"-f"`

	// -g
	AllowRemoteConnectionsToLocalForwardedPorts bool `flag:"-g"`

	// -I <pkcs11>
	PKCS11 string `flag:"-I"`

	// -i <identity>
	Identity string `flag:"-i"`

	// -K
	GSSAPI bool `flag:"-K"`

	// -k
	NoGSSAPI bool `flag:"-k"`

	// -L [<bind_address>:]<port>:<host>:<hostport>
	LocalForward string `flag:"-L"` // FIXME data structure

	// -l <login_name>
	Login string `flag:"-l"`

	// -M
	Master bool `flag:"-M"`

	// -m <mac_spec>
	MACSpec string `flag:"-m"`

	// -N
	NoRemoteCommand bool `flag:"-N"`

	// -n
	NoReadStdin bool `flag:"-n"`

	// -O <ctl_cmd>
	ControlCommand string `flag:"-O"`

	// -o <option>
	Options SSHOptions `flag:"-"`

	// -p <port>
	Port int `flag:"-p"`

	// -q
	Quiet bool `flag:"-q"`

	// -R [<bind_address>:]<port>:<host>:<hostport>
	RemoteForward string `flag:"-R"` // FIXME data structure

	// -S <ctl_path>
	ControlPath string `flag:"-S"`

	// -s
	Subsystem bool `flag:"-s"`

	// -T
	NoTTY bool `flag:"-T"`

	// -t
	TTY bool `flag:"-t"`

	// -v, -vv, and -vvv
	Verbose  bool `flag:"-v"`
	Verbose1 bool `flag:"-v"`
	Verbose2 bool `flag:"-vv"`
	Verbose3 bool `flag:"-vvv"`

	// -W <host>:<port>
	ForwardStdinStdout string `flag:"-W"` // FIXME data structure

	// -w <local_tun>:<remote_tun>
	Tunnel string `flag:"-w"` // FIXME data structure

	// -X
	X11 bool `flag:"-X"`

	// -x
	NoX11 bool `flag:"-x"`

	// -Y
	TrustedX11 bool `flag:"-Y"`

	// -y
	Syslog bool `flag:"-y"`

	// [<username>@]<hostname>
	Hostname string `pos:"last"`

	// <command>
	Command []string `pos:"last"`
}

// SSHOptions is a map of options as specified in ssh_config(5) files.
type SSHOptions map[string]string

func (m SSHOptions) String() string {
	options := make(sort.StringSlice, 0, len(m))
	for k, v := range m {
		options = append(options, fmt.Sprintf("%s=%s", k, v))
	}
	options.Sort()
	args := make([]string, 2*len(m))
	for i, option := range options {
		args[2*i] = "-o"
		args[2*i+1] = option
	}
	return strings.Join(args, " ")
}
