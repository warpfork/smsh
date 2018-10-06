package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

func main() {
	ctx := context.Background()
	if err := Main(ctx, os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		switch e2 := err.(type) {
		case ErrChildExit:
			fmt.Fprintf(os.Stderr, "smsh: %s", err)
			if e2.Signal != 0 {
				os.Exit(128 + e2.Signal)
			}
			os.Exit(e2.Code)
		case ErrInternal:
			fmt.Fprintf(os.Stderr, "smsh: %s", err)
			os.Exit(125)
		default:
			panic(err)
		}
	}
}

func Main(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	cmdStrm := bytes.NewBufferString(strings.Join(args, "\n"))
	runner, _ := interp.New(interp.StdIO(stdin, stdout, stderr))
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
		return true
	}
	return parser.Stmts(cmdStrm, fn)
}
