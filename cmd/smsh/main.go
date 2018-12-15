package main

import (
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
	runner, _ := interp.New(
		interp.StdIO(stdin, stdout, stderr),
		interp.Module(execTool),
		interp.Params("-e"), // TODO this doesn't do anything?
		interp.Params("-u"), // TODO check if this does either
	)

	parser := syntax.NewParser()
	for _, arg := range args {
		file, err := parser.Parse(strings.NewReader(arg), "")
		if err != nil {
			return err
		}
		for _, stmt := range file.Stmts {
			//fmt.Printf(":: %#v\n", stmt)
			if err := runner.Run(ctx, stmt); err != nil {
				switch err.(type) {
				case ErrInternal, ErrChildExit:
					return err
				default:
					return ErrInternal{err}
				}
			}
		}
	}
	return nil
}
