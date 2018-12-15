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
			fmt.Fprintf(os.Stderr, "smsh: %s\n", err)
			if e2.Signal != 0 {
				os.Exit(128 + e2.Signal)
			}
			os.Exit(e2.Code)
		case ErrInternal:
			fmt.Fprintf(os.Stderr, "smsh: %s\n", err)
			os.Exit(125)
		default:
			panic(err)
		}
	}
}

func Main(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) (halt error) {
	return Run(ctx, args[1:], stdin, stdout, stderr)
}

func Run(ctx context.Context, cmdStrs []string, stdin io.Reader, stdout, stderr io.Writer) (halt error) {
	parser := syntax.NewParser()
	for _, cmdStr := range cmdStrs {
		file, err := parser.Parse(strings.NewReader(cmdStr), "")
		if err != nil {
			return err
		}
		err = RunOne(ctx, file, stdin, stdout, stderr)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunOne(ctx context.Context, cmd *syntax.File, stdin io.Reader, stdout, stderr io.Writer) (halt error) {
	runner, _ := interp.New(
		interp.StdIO(stdin, stdout, stderr),
		interp.Module(execTool),
		interp.Params("-e"), // TODO this doesn't do anything?
		interp.Params("-u"), // TODO check if this does either
	)
	for _, stmt := range cmd.Stmts {
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
	return nil
}
