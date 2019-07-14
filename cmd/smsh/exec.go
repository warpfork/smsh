package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"mvdan.cc/sh/expand"
	"mvdan.cc/sh/interp"
)

var execTool = interp.ModuleExec(func(ctx context.Context, path string, args []string) error {
	mc, _ := interp.FromModuleContext(ctx)
	if path == "" {
		fmt.Fprintf(mc.Stderr, "%q: executable file not found in $PATH\n", args[0])
		return ErrChildExit{args[0], 127, 0}
	}
	cmd := exec.Cmd{
		Path:   path,
		Args:   args,
		Env:    execEnv(mc.Env),
		Dir:    mc.Dir,
		Stdin:  mc.Stdin,
		Stdout: mc.Stdout,
		Stderr: mc.Stderr,
	}

	err := cmd.Start()
	if err == nil {
		if done := ctx.Done(); done != nil {
			go func() {
				<-done

				if mc.KillTimeout <= 0 || runtime.GOOS == "windows" {
					_ = cmd.Process.Signal(os.Kill)
					return
				}

				// TODO: don't temporarily leak this goroutine
				// if the program stops itself with the
				// interrupt.
				go func() {
					time.Sleep(mc.KillTimeout)
					_ = cmd.Process.Signal(os.Kill)
				}()
				_ = cmd.Process.Signal(os.Interrupt)
			}()
		}

		err = cmd.Wait()
	}

	switch x := err.(type) {
	case *exec.ExitError: // started, but errored
		if status, ok := x.Sys().(syscall.WaitStatus); ok {
			if status.Signaled() {
				return ErrChildExit{args[0], 0, int(status.Signal())}
			}
			return ErrChildExit{args[0], status.ExitStatus(), 0}
		}
		// default to 1 if the OS doens't have exit statuses
		return ErrChildExit{args[0], 1, 0}
	case *exec.Error: // did not start
		return ErrInternal{err}
	case nil:
		return nil
	default:
		return ErrInternal{err}
	}
})

func execEnv(env expand.Environ) []string {
	list := make([]string, 32)
	env.Each(func(name string, vr expand.Variable) bool {
		list = append(list, name+"="+vr.String())
		return true
	})
	return list
}
