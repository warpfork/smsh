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
	err := Main(context.Background(), args, stdin, stdout, stderr)
	return stdout.String(), stderr.String(), err
}

func Test(t *testing.T) {
	stdout, stderr, err := runMain("echo foo", "echo bar", "echo baz")
	Wish(t, stdout, ShouldEqual, Dedent(`
		foo
		bar
		baz
	`))
	Wish(t, stderr, ShouldEqual, "")
	Wish(t, err, ShouldEqual, nil)
}
