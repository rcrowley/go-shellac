package shellac

import "testing"

func TestArgs(t *testing.T) {
	testArgs(t, []string{}, Args(test{}))
}

func TestArgsFlagArray(t *testing.T) {
	testArgs(t, []string{"-flag-array", "hi", "hi"}, Args(test{
		FlagArray: [2]string{"hi", "hi"},
	}))
}

func TestArgsFlag(t *testing.T) {
	testArgs(t, []string{"-flag", "hi"}, Args(test{Flag: "hi"}))
}

func TestArgsFlagBool(t *testing.T) {
	testArgs(t, []string{"-flag-bool"}, Args(test{FlagBool: true}))
}

func TestArgsFlagEmptySep(t *testing.T) {
	testArgs(t, []string{"-fhi"}, Args(test{FlagEmptySep: "hi"}))
}

func TestArgsFlagInt(t *testing.T) {
	testArgs(t, []string{"-flag-int", "47"}, Args(test{FlagInt: 47}))
}

func TestArgsFlagSep(t *testing.T) {
	testArgs(t, []string{"-flag-sep=hi"}, Args(test{FlagSep: "hi"}))
}

func TestArgsFlagSlice(t *testing.T) {
	testArgs(t, []string{"-flag-slice", "hi", "hi"}, Args(test{
		FlagSlice: []string{"hi", "hi"},
	}))
}

func TestArgsFlagZero(t *testing.T) {
	testArgs(t, []string{}, Args(test{Flag: ""}))
	testArgs(t, []string{}, Args(test{FlagInt: 0}))
}

func TestArgsFlagZeroPtr(t *testing.T) {
	testArgs(t, []string{"-flag-int-ptr", "0"}, Args(test{
		FlagIntPtr: NewInt(0)},
	))
	testArgs(t, []string{"-flag-ptr", ""}, Args(test{FlagPtr: NewString("")}))
}

func TestArgsPosFirst(t *testing.T) {
	testArgs(t, []string{"first", "-flag", "hi"}, Args(test{
		Flag:     "hi",
		PosFirst: "first"},
	))
}

func TestArgsPosFirstInt(t *testing.T) {
	testArgs(t, []string{"47", "-flag", "hi"}, Args(test{
		Flag:        "hi",
		PosFirstInt: 47,
	}))
}

func TestArgsPosFirstZero(t *testing.T) {
	testArgs(t, []string{"-flag", "hi"}, Args(test{
		Flag:     "hi",
		PosFirst: ""},
	))
	testArgs(t, []string{"-flag", "hi"}, Args(test{
		Flag:        "hi",
		PosFirstInt: 0,
	}))
}

func TestArgsPosFirstZeroPtr(t *testing.T) {
	testArgs(t, []string{"0", "-flag", "hi"}, Args(test{
		Flag:           "hi",
		PosFirstIntPtr: NewInt(0)},
	))
	testArgs(t, []string{"", "-flag", "hi"}, Args(test{
		Flag:        "hi",
		PosFirstPtr: NewString(""),
	}))
}

func TestArgsPosLast(t *testing.T) {
	testArgs(t, []string{"-flag", "hi", "last"}, Args(test{
		Flag:    "hi",
		PosLast: "last",
	}))
}

func TestArgsPosLastInt(t *testing.T) {
	testArgs(t, []string{"-flag", "hi", "47"}, Args(test{
		Flag:       "hi",
		PosLastInt: 47,
	}))
}

func TestArgsPosLastZero(t *testing.T) {
	testArgs(t, []string{"-flag", "hi"}, Args(test{
		Flag:    "hi",
		PosLast: "",
	}))
	testArgs(t, []string{"-flag", "hi"}, Args(test{
		Flag:       "hi",
		PosLastInt: 0,
	}))
}

func TestArgsPosLastZeroPtr(t *testing.T) {
	testArgs(t, []string{"-flag", "hi", "0"}, Args(test{
		Flag:          "hi",
		PosLastIntPtr: NewInt(0),
	}))
	testArgs(t, []string{"-flag", "hi", ""}, Args(test{
		Flag:       "hi",
		PosLastPtr: NewString(""),
	}))
}

func TestArgsPosWrong(t *testing.T) {
	testArgs(t, []string{"-flag", "hi"}, Args(test{
		Flag:     "hi",
		PosWrong: "wrong",
	}))
}

func TestCommand(t *testing.T) {
	cmd := Command(test{})
	if "/usr/bin/test" != cmd.Path || "test" != cmd.Args[0] {
		t.Fatal(cmd)
	}
}

func TestCommandDefault(t *testing.T) {
	cmd := Command(testDefault{})
	if "testdefault" != cmd.Path || "testdefault" != cmd.Args[0] {
		t.Fatal(cmd)
	}
}

func TestSudoCommand(t *testing.T) {
	cmd := Command(test{})
	cmd.Sudo()
	if "/usr/bin/sudo" != cmd.Path || "sudo" != cmd.Args[0] {
		t.Fatal(cmd)
	}
	if "test" != cmd.Args[1] {
		t.Fatal(cmd)
	}
}

type test struct {
	_              struct{}  `command:"test"`
	Flag           string    `flag:"-flag"`
	FlagArray      [2]string `flag:"-flag-array"`
	FlagBool       bool      `flag:"-flag-bool"`
	FlagEmptySep   string    `flag:"-f" sep:"-"`
	FlagInt        int       `flag:"-flag-int"`
	FlagIntPtr     *int      `flag:"-flag-int-ptr"`
	FlagPtr        *string   `flag:"-flag-ptr"`
	FlagSep        string    `flag:"-flag-sep" sep:"="`
	FlagSlice      []string  `flag:"-flag-slice"`
	PosFirst       string    `pos:"first"`
	PosFirstInt    int       `pos:"first"`
	PosFirstIntPtr *int      `pos:"first"`
	PosFirstPtr    *string   `pos:"first"`
	PosLast        string    `pos:"last"`
	PosLastInt     int       `pos:"last"`
	PosLastIntPtr  *int      `pos:"last"`
	PosLastPtr     *string   `pos:"last"`
	PosWrong       string    `pos:"wrong"`
}

type testDefault struct{}

func testArgs(t *testing.T, expected, actual []string) {
	if len(expected) != len(actual) {
		t.Fatal(expected, actual)
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			t.Fatal(expected, actual)
		}
	}
	t.Log(actual)
}
