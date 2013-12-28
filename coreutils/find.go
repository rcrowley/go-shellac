package coreutils

import "fmt"

// find(1)
//
// There is no support for the complex logical expressions that are possible
// with find(1).  If you need to execute such commands, use the exec package.
type Find struct {

	// The three basic modes of operation on symbolic links.  The last of these
	// flags to be specified to find(1) wins and nullifies any others; because
	// their order here is fixed, consider them mutually exclusive.
	//
	// Note also that -follow is not supported because it's incompatible with
	// our need to specify arguments in a predictable order.
	DoNotFollowSymlinks   bool `flag:"-P" pos:"first"`
	FollowSymlinks        bool `flag:"-L" pos:"first"`
	FollowInitialSymlinks bool `flag:"-H" pos:"first"`

	// -D <debugoptions>
	DebugOptions string `flag:"-D" pos:"first"`

	// -O<level>
	Optimization int `flag:"-O" pos:"first" sep:"-"`

	// The list of directories from which find(1) will begin its traversal.
	Dirnames []string `pos:"first"`

	// -daystart
	DayStart bool `flag:"-daystart"`

	// -depth
	DepthFirst bool `flag:"-depth"`

	// -ignore_readdir_race and -noignore_readdir_race
	IgnoreReaddirRace bool `flag:"-ignore_readdir_race"`
	//NoIgnoreReaddirRace bool `flag:"-noignore_readdir_race"`

	// -maxdepth <levels>
	MaxDepth *int `flag:"-maxdepth"`

	// -mindepth <levels>
	MinDepth *int `flag:"-mindepth"`

	// -noleaf
	NoLeaf bool `flag:"-noleaf"`

	// -regextype <type>
	RegexType string `flag:"-regextype"`

	// -warn and -nowarn
	Warn bool `flag:"-warn"`
	// NoWarn bool `flag:"-nowarn"`

	// -xdev (formerly known as -mount)
	XDev bool `flag:"-xdev"`

	// -amin <n>
	AccessedMinutesAgo *FindN `flag:"-amin"`

	// -anewer <file>
	AccessedSinceFile string `flag:"-anewer"`

	// -atime <n>
	AccessedDaysAgo *FindN `flag:"-atime"`

	// -cmin <n>
	ChangedMinutesAgo *FindN `flag:"-cmin"`

	// -cnewer <file>
	ChangedSinceFile string `flag:"-cnewer"`

	// -ctime <n>
	ChangedDaysAgo *FindN `flag:"-ctime"`

	// -empty
	Empty bool `flag:"-empty"`

	// -executable
	Executable bool `flag:"-executable"`

	// -false
	False bool `flag:"-false"`

	// -fstype <type>
	FilesystemType string `flag:"-fstype"`

	// -gid <n>
	GID *FindN `flag:"-gid"`

	// -group <gname>
	Group string `flag:"-group"`

	// -ilname <pattern>
	SymlinkTargetCaseInsensitive string `flag:"-ilname"`

	// -iname <pattern>
	NameCaseInsensitive string `flag:"-iname"`

	// -inum <n>
	Inode *FindN `flag:"-inum"`

	// -iregex <pattern>
	RegexCaseInsensitive string `flag:"-iregex"`

	// -iwholename <pattern>
	WholenameCaseInsensitive string `flag:"-iwholename"`

	// -links <n>
	Links *FindN `flag:"-links"`

	// -lname
	LinkName string `flag:"-lname"`

	// -mmin <n>
	ModifiedMinutesAgo *FindN `flag:"-mmin"`

	// -mtime <n>
	ModifiedDaysAgo *FindN `flag:"-mtime"`

	// -name <pattern>
	Name string `flag:"-name"`

	// -newer <file>
	ModifiedSinceFile string `flag:"-newer"`
	Newer             string `flag:"-newer"`

	// TODO -newerXY <reference>

	// -nogroup
	UnnamedGroup bool `flag:"-nogroup"`

	// -nouser
	UnnamedUser bool `flag:"-nouser"`

	// -path <pattern>
	Path string `flag:"-path"`

	// -perm <mode>
	Mode *int `flag:"-perm" format:"%o"`

	// -perm -<mode>
	ModeMaskAll *int `flag:"-perm" format:"-%o"`

	// -perm /<mode>
	ModeMaskAny *int `flag:"-perm" format:"/%o"`

	// -readable
	Readable bool `flag:"-readable"`

	// -regex <pattern>
	Regex string `flag:"-regex"`

	// -samefile <name>
	SameFile string `flag:"-samefile"`

	// -size <n>[cwbkMG] (but only byte sizes are supported)
	// TODO support more than just the c (bytes) suffix.
	Size *FindN `flag:"-size" format:"%sc"`

	// -true
	True bool `flag:"-true"`

	// -type <c>
	Type FindType `flag:"-type"`

	// -uid <n>
	UID *FindN `flag:"-uid"`

	// -used <n>
	Used *FindN `flag:"-used"`

	// -user <uname>
	User string `flag:"-user"`

	// -writable
	Writable bool `flag:"-writable"`

	// -xtype <c>
	XType FindType `flag:"-xtype"`

	// -delete
	Delete bool `flag:"-delete" pos:"last"`

	// -exec <command> ; or -exec <command> +
	Exec []string `flag:"-exec" pos:"last"`

	// -execdir <command> ; or -execdir <command> +
	ExecDir []string `flag:"-execdir" pos:"last"`

	// -fls <file>
	Fls string `flag:"-fls" pos:"last"`

	// -fprint <file>
	Fprint string `flag:"-fprint" pos:"last"`

	// -fprint0 <file>
	Fprint0 string `flag:"-fprint0" pos:"last"`

	// -fprintf <file> <format>
	Fprintf [2]string `flag:"-fprintf" pos:"last"`

	// -ls
	Ls bool `flag:"-ls" pos:"last"`

	// -ok <command> ;
	// -ok <command> ; or -ok <command> +
	OK []string `flag:"-ok" pos:"last"`

	// -okdir <command> ; or -okdir <command> +
	OKDir []string `flag:"-okdir" pos:"last"`

	// -print
	Print bool `flag:"-print" pos:"last"`

	// -print0
	Print0 bool `flag:"-print0" pos:"last"`

	// -printf <format>
	Printf string `flag:"-printf" pos:"last"`

	// -prune
	Prune bool `flag:"-prune" pos:"last"`

	// -quit
	Quit bool `flag:"-quit" pos:"last"`
}

// FindExec is a slice of strings representing the arguments to find(1)'s
// -exec, -execdir, -ok, and -okdir options, exposed as Exec, ExecDir, OK, and
// OKDir in Find.
type FindExec []string

// NewFindExec constructs a FindExec slice representing the arguments to
// find(1)'s -exec, -execdir, -ok, and -okdir options, exposed as Exec,
// ExecDir, OK, and OKDir in Find.
func NewFindExec(mode FindExecMode, args ...string) FindExec {
	return FindExec(append(args, string(mode)))
}

// FindExecMode is an enumeration of all possible qualifiers to find(1)'s
// -exec, -execdir, -ok, and -okdir options, exposed as Exec, ExecDir, OK, and
// OKDir in Find.
type FindExecMode string

var (
	FindExecOne  FindExecMode = ";"
	FindExecMany FindExecMode = "+"
)

// FindN is a string, conventionally of the form "n", "+n", or "-n" where n is
// an integer.
type FindN struct {
	Mode FindNMode
	N    int
}

// NewFindN constructs a FindN string in the conventional form.
func NewFindN(mode FindNMode, n int) *FindN {
	return &FindN{mode, n}
}

func (n FindN) String() string {
	return fmt.Sprintf("%s%d", n.Mode, n.N)
}

// FindNMode is an enumeration of all possible qualifiers to find(1)'s <n>
// arguments.
type FindNMode string

var (
	FindExact       FindNMode = ""
	FindGreaterThan FindNMode = "+"
	FindLessThan    FindNMode = "-"
)

// FindType is an enumeration of all possible values of find(1)'s -type and
// -xtype options, exposed as Type and XType in Find.
type FindType string

var (
	FindBlock     FindType = "b"
	FindCharacter FindType = "c"
	FindDirectory FindType = "d"
	FindPipe      FindType = "p"
	FindFile      FindType = "f"
	FindLink      FindType = "l"
	FindSocket    FindType = "s"
	FindDoor      FindType = "D"
)
