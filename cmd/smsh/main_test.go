package main

import (
	"bytes"
	"context"
	"testing"

	. "github.com/warpfork/go-wish"
)

func runMain(args ...string) (string, string, error) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	err := Run(context.Background(), args, stdin, stdout, stderr)
	return stdout.String(), stderr.String(), err
}

func TestHappyEcho(t *testing.T) {
	stdout, stderr, err := runMain("echo foo", "echo bar", "echo baz")
	Wish(t, stdout, ShouldEqual, Dedent(`
		foo
		bar
		baz
	`))
	Wish(t, stderr, ShouldEqual, "")
	Wish(t, err, ShouldEqual, nil)
}

func TestHappyEcho2(t *testing.T) {
	stdout, stderr, err := runMain("echo foo; echo bar", "echo baz")
	Wish(t, stdout, ShouldEqual, Dedent(`
		foo
		bar
		baz
	`))
	Wish(t, stderr, ShouldEqual, "")
	Wish(t, err, ShouldEqual, nil)
}

func TestExitOnError(t *testing.T) {
	stdout, stderr, err := runMain("echo foo", "thisshouldnotbeacommand", "echo baz")
	Wish(t, stdout, ShouldEqual, Dedent(`
		foo
	`))
	Wish(t, stderr, ShouldEqual, Dedent(`
		"thisshouldnotbeacommand": executable file not found in $PATH
	`))
	Wish(t, err, ShouldEqual, ErrChildExit{"thisshouldnotbeacommand", 127, 0})
}

func TestPipes(t *testing.T) {
	stdout, stderr, err := runMain("echo foo | cat -")
	Wish(t, stdout, ShouldEqual, Dedent(`
		foo
	`))
	Wish(t, stderr, ShouldEqual, "")
	Wish(t, err, ShouldEqual, nil)
}

func TestComments(t *testing.T) {
	t.Run("single end of line comment", func(t *testing.T) {
		stdout, stderr, err := runMain("echo foo | cat - # cmnt")
		Wish(t, stdout, ShouldEqual, "foo\n")
		Wish(t, stderr, ShouldEqual, "")
		Wish(t, err, ShouldEqual, nil)
	})
}

func TestFuncsAndSubshells(t *testing.T) {
	t.Run("single end of line comment", func(t *testing.T) {
		stdout, stderr, err := runMain(`(function foo { :; echo fa; }; foo)`)
		Wish(t, stdout, ShouldEqual, "fa\n")
		Wish(t, stderr, ShouldEqual, "")
		Wish(t, err, ShouldEqual, nil)
	})
}
