package shellac

import (
	"github.com/rcrowley/go-shellac/coreutils"
	"testing"
)

func TestFind(t *testing.T) {
	testArgs(t, []string{}, Args(coreutils.Find{}))
	testArgs(t, []string{"."}, Args(coreutils.Find{Dirnames: []string{"."}}))
}

func TestFindExec(t *testing.T) {
	testArgs(t, []string{"-exec", "cat", "{}", ";"}, Args(coreutils.Find{
		Exec: coreutils.NewFindExec(coreutils.FindExecOne, "cat", "{}"),
	}))
	testArgs(t, []string{"-exec", "grep", "foo bar", "{}", "+"}, Args(
		coreutils.Find{
			Exec: coreutils.NewFindExec(
				coreutils.FindExecMany,
				"grep",
				"foo bar",
				"{}",
			),
		},
	))
}

func TestFindMode(t *testing.T) {
	testArgs(t, []string{"-perm", "644"}, Args(coreutils.Find{
		Mode: NewInt(0644),
	}))
	testArgs(t, []string{"-perm", "-644"}, Args(coreutils.Find{
		ModeMaskAll: NewInt(0644),
	}))
	testArgs(t, []string{"-perm", "/644"}, Args(coreutils.Find{
		ModeMaskAny: NewInt(0644),
	}))
}

func TestFindN(t *testing.T) {
	testArgs(t, []string{"-gid", "0"}, Args(coreutils.Find{
		GID: coreutils.NewFindN(coreutils.FindExact, 0),
	}))
	testArgs(t, []string{"-links", "+3"}, Args(coreutils.Find{
		Links: coreutils.NewFindN(coreutils.FindGreaterThan, 3),
	}))
	testArgs(t, []string{"-uid", "-1000"}, Args(coreutils.Find{
		UID: coreutils.NewFindN(coreutils.FindLessThan, 1000),
	}))
}

func TestFindSymlinks(t *testing.T) {
	testArgs(t, []string{"-P"}, Args(coreutils.Find{
		DoNotFollowSymlinks: true,
	}))
	testArgs(t, []string{"-L"}, Args(
		coreutils.Find{FollowSymlinks: true}))
	testArgs(t, []string{"-H"}, Args(coreutils.Find{
		FollowInitialSymlinks: true,
	}))
}
