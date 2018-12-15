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
	cmdAll := strings.Join(args, "\n")
	cmdStrm := bytes.NewBufferString(cmdAll)
	runner, _ := interp.New(
		interp.StdIO(stdin, stdout, stderr),
		interp.Module(execTool),
		interp.Params("-e"), // TODO this doesn't do anything?
		interp.Params("-u"), // TODO check if this does either
	)
	fn := func(s *syntax.Stmt) bool {
		//fmt.Printf(":: %#v\n", cmdAll[s.Pos().Offset():s.End().Offset()])
		if err := runner.Run(ctx, s); err != nil {
			switch err.(type) {
			case ErrInternal, ErrChildExit:
				halt = err
			default:
				halt = ErrInternal{err}
			}
			return false
		}
		return true
	}
	parser := syntax.NewParser()
	if err := parser.Stmts(cmdStrm, fn); err != nil {
		return err
	}
	return
}
