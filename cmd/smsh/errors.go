package main

import (
	"fmt"
)

var (
	_ error = ErrChildExit{}
	_ error = ErrInternal{}
)

type ErrChildExit struct {
	Name   string
	Code   int
	Signal int
}

func (e ErrChildExit) Error() string {
	if e.Signal != 0 {
		return fmt.Sprintf("process %q halted by signal %d", e.Name, e.Signal)
	}
	return fmt.Sprintf("process %q exited with code %d", e.Name, e.Code)
}

type ErrInternal struct {
	error
}
