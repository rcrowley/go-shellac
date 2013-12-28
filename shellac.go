// Package shellac provides a declarative, strongly-typed API for executing
// shell commands.
package shellac

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

// Args returns a slice of strings of the arguments to the command described by
// the given interface value, which should be a struct or pointer to a struct.
//
// Zero-valued fields are omitted.  To have zero values like 0 or "" included,
// use a pointer type and the NewInt or NewString helpers or define helpers of
// your own.
//
// Arrays and slices of strings are used as-is.  Scalar-valued fields are made
// to be strings by the fmt package.  If the field has a format tag, that is
// used as the format argument to fmt.Sprintf; otherwise the standard %v format
// is used.
//
// Fields with a flag tag are stringified as above.  The value of the flag tag
// is used as a prefix unless the value of the flag tag is -.
//
// Boolean flags return the flag tag itself if the field is true and the empty
// slice otherwise.
//
// All other flags include a separator, as configured by the sep tag, between
// the flag and the field value.  By default they're separated by a single
// space; the tag "-" removes the separator; any other value is used literally.
//
// See <https://github.com/rcrowley/go-shellac/blob/master/shellac_test.go> for
// examples.
func Args(i interface{}) []string {
	v := reflect.ValueOf(i)
	if reflect.Ptr == v.Kind() {
		v = v.Elem()
	}
	t := v.Type()
	fields := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if pos := tag.Get("pos"); "first" == pos {
			fields = append(fields, field(t.Field(i), v.Field(i))...)
		}
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if "_" == f.Name {
			continue
		}
		tag := f.Tag
		if "" == tag.Get("flag") {
			continue
		}
		if pos := tag.Get("pos"); "first" != pos && "last" != pos {
			fields = append(fields, field(t.Field(i), v.Field(i))...)
		}
	}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if pos := tag.Get("pos"); "last" == pos {
			fields = append(fields, field(t.Field(i), v.Field(i))...)
		}
	}
	return fields
}

// Cmd wraps exec.Cmd to add convenience methods.
type Cmd struct {
	exec.Cmd
}

// Command returns a *Cmd (with standard input, output, and error connected)
// as described by the given interface value, which should be a struct or
// pointer to a struct.
func Command(i interface{}) *Cmd {
	t := reflect.TypeOf(i)
	if reflect.Ptr == t.Kind() {
		t = t.Elem()
	}
	cmd := &Cmd{*exec.Command(command(t), Args(i)...)}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// ChannelStdin connects standard input to a channel.
func (cmd *Cmd) ChannelStdin(stdin <-chan string) {
	cmd.Stdin = NewChanReader(stdin)
}

// ChannelStdout connects standard output to a channel.
func (cmd *Cmd) ChannelStdout(stdout chan<- string) {
	cmd.Stdout = NewChanWriter(stdout)
}

// ChannelStderr connects standard errput to a channel.
func (cmd *Cmd) ChannelStderr(stderr chan<- string) {
	cmd.Stderr = NewChanWriter(stderr)
}

// Log logs the command (bolded if standard error is a TTY) to standard error.
// This is sort of like what make(1) or sh(1) with -x do.
func (cmd *Cmd) Log() {
	fi, err := os.Stderr.Stat()
	if nil != err {
		panic(err)
	}
	var format string
	if 0 == fi.Mode()&os.ModeCharDevice {
		format = "%s\n"
	} else {
		format = "\033[1m%s\033[0m\n"
	}
	fmt.Fprintf(os.Stderr, format, strings.Join(cmd.Args, " "))
}

// Run logs and runs a shell command.
func (cmd *Cmd) Run() error {
	cmd.Log()
	defer cmd.closeStdoutStderr()
	return cmd.Cmd.Run()
}

// Sudo modifies a shell command to be run as root via sudo(8).
func (cmd *Cmd) Sudo() {
	sudo, err := exec.LookPath("sudo")
	if nil != err {
		panic(err)
	}
	args := make([]string, 1+len(cmd.Args))
	args[0] = "sudo"
	copy(args[1:], cmd.Args)
	cmd.Args = args
	cmd.Path = sudo
}

// closeStdoutStderr calls Close on either or both of standard output and error
// that is using a ChanWriter.
func (cmd *Cmd) closeStdoutStderr() error {
	if w, ok := cmd.Stdout.(*ChanWriter); ok {
		if err := w.Close(); nil != err {
			return err
		}
	}
	if w, ok := cmd.Stderr.(*ChanWriter); ok {
		if err := w.Close(); nil != err {
			return err
		}
	}
	return nil
}

// NewInt returns a pointer to the given integer.
func NewInt(i int) *int {
	return &i
}

// NewString returns a pointer to the given string.
func NewString(s string) *string {
	return &s
}

// Run constructs and runs a shell command from the given interface value or
// simply runs a Cmd or exec.Cmd.
func Run(i interface{}) error {
	if cmd, ok := i.(*Cmd); ok {
		return cmd.Run()
	}
	if execCmd, ok := i.(*exec.Cmd); ok {
		cmd := &Cmd{*execCmd}
		return cmd.Run()
	}
	return Command(i).Run()
}

// Sudo constructs and runs a shell command from the given interface value or
// simply runs a Cmd or exec.Cmd.  In any case, the command is run as root via
// sudo(8).
func Sudo(i interface{}) error {
	if cmd, ok := i.(*Cmd); ok {
		cmd.Sudo()
		return cmd.Run()
	}
	if execCmd, ok := i.(*exec.Cmd); ok {
		cmd := &Cmd{*execCmd}
		cmd.Sudo()
		return cmd.Run()
	}
	cmd := Command(i)
	cmd.Sudo()
	return cmd.Run()
}

// command returns the name of the command either from the command tag on any
// field or from the lowercase name of the struct type.
func command(t reflect.Type) string {
	for i := 0; i < t.NumField(); i++ {
		if command := t.Field(i).Tag.Get("command"); "" != command {
			return command
		}
	}
	return strings.ToLower(t.Name())
}

// field returns a slice of strings representing the reflect.StructField f from
// reflect.Value v as it should be written in a shell command.
func field(f reflect.StructField, v reflect.Value) []string {
	if "" != f.PkgPath {
		return []string{}
	}
	t := f.Type
	k := t.Kind()
	switch k {
	case reflect.Ptr:
		if v.IsNil() {
			return []string{}
		}
		v = v.Elem()
	case reflect.Chan, reflect.Map, reflect.Slice:
		if v.IsNil() {
			return []string{}
		}
	case reflect.Struct:
		if 0 == t.NumField() {
			return []string{}
		}
	default:
		if reflect.Zero(t).Interface() == v.Interface() {
			return []string{}
		}
	}
	flag := f.Tag.Get("flag")
	if "" != flag && reflect.Bool == k && v.Bool() {
		return []string{flag}
	}
	args := []string{}
	if "" != flag {
		args = append(args, flag)
	}
	if reflect.Array == k && reflect.String == t.Elem().Kind() {
		for i := 0; i < v.Len(); i++ {
			args = append(args, v.Index(i).String())
		}
		return args
	}
	if reflect.Slice == k && reflect.String == t.Elem().Kind() {
		return append(args, v.Interface().([]string)...)
	}
	format := f.Tag.Get("format")
	if "" == format {
		format = "%v"
	}
	arg := fmt.Sprintf(format, v.Interface())
	if "" == flag || "-" == flag {
		return []string{arg}
	}
	switch sep := f.Tag.Get("sep"); sep {
	case "":
		return []string{flag, arg}
	case "-":
		return []string{fmt.Sprintf("%s%s", flag, arg)}
	default:
		return []string{fmt.Sprintf("%s%s%s", flag, sep, arg)}
	}
}
