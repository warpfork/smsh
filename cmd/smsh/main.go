package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

func main() {
	if err := zoom(); err != nil {
		fmt.Fprintf(os.Stderr, "smsh: %s", err)
		os.Exit(1)
	}
}

func zoom() error {
	ctx := context.Background()
	pr := &promptReader{os.Stdin, os.Stdout, true}
	runner, _ := interp.New(interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	parser := syntax.NewParser()
	fn := func(s *syntax.Stmt) bool {
		if err := runner.Run(ctx, s); err != nil {
			switch x := err.(type) {
			case interp.ShellExitStatus:
				os.Exit(int(x))
			case interp.ExitStatus:
			default:
				fmt.Fprintln(runner.Stderr, err)
				os.Exit(1)
			}
		}
		pr.first = true
		return true
	}
	return parser.Stmts(pr, fn)
}

type promptReader struct {
	io.Reader
	io.Writer
	first bool
}

func (pr *promptReader) Read(p []byte) (int, error) {
	if pr.first {
		fmt.Fprintf(pr.Writer, "$ ")
		pr.first = false
	} else {
		fmt.Fprintf(pr.Writer, "> ")
	}
	return pr.Reader.Read(p)
}
