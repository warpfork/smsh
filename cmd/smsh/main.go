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

func Main(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) (halt error) {
	// TODO : I think we're going to have to feed things piecewise if we want to produce reports in the same granularity as the args.
	//  The parser is clever enough to split statements, and that's... kind of not exactly what want here.
	cmdStrm := bytes.NewBufferString(strings.Join(args, "\n"))
	runner, _ := interp.New(
		interp.StdIO(stdin, stdout, stderr),
		interp.Params("-e"), // TODO this doesn't do anything?
		interp.Params("-u"), // TODO check if this does either
	)
	fn := func(s *syntax.Stmt) bool {
		fmt.Printf(":: %#v\n", s)
		if err := runner.Run(ctx, s); err != nil {
			nodeStr := "todo:restring-the-cmd" // TODO this is hard to string back up!  fmt.Sprintf("%#v", s.Cmd) isn't even close.  offsets might do it?  uff.
			switch x := err.(type) {
			case interp.ShellExitStatus:
				// TODO it's not clear to me why this and ExitStatus are distinct,
				//  so for now I'm just boxing them both into the same thing.
				halt = ErrChildExit{nodeStr, int(x), 0}
				return false
			case interp.ExitStatus:
				// TODO it looks like interp.DefaultExec doesn't really handle signals
				//  as completely as it could... we should consider some PRs to that.
				halt = ErrChildExit{nodeStr, int(x), 0}
				return false
			default:
				halt = ErrInternal{err}
				return false
			}
		}
		return true
	}
	parser := syntax.NewParser()
	if err := parser.Stmts(cmdStrm, fn); err != nil {
		return err
	}
	return
}
