package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/subosito/gotenv"
)

func main() {
	path := flag.String("f", "", "filepath to .env")
	flag.Parse()

	env := parseEnv(*path)
	for key, value := range env {
		os.Setenv(key, value)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("nothing to start")
		os.Exit(0)
	}

	name := args[0]
	var subArgs []string
	if len(args) > 1 {
		subArgs = args[1:]
	}

	cmd := exec.Command(name, subArgs...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh)

	waitCh := make(chan error)
	go func() {
		waitCh <- cmd.Wait()

		close(waitCh)
	}()

	for {
		select {
		case s := <-signCh:
			err = cmd.Process.Signal(s)
			if err != nil {
				fmt.Println("error sendind signal", s, err)
			}
		case err = <-waitCh:
			exitErr, ok := err.(*exec.ExitError)
			if ok {
				childStatus := exitErr.Sys().(syscall.WaitStatus)

				os.Exit(childStatus.ExitStatus())
			} else {
				fmt.Println(err.Error())

				os.Exit(1)
			}
		}
	}
}

func parseEnv(path string) map[string]string {
	if path == "" {
		return nil
	}

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	// beacuse of we only read a file
	defer f.Close() //nolint:gosec

	env, err := gotenv.StrictParse(f)
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	return env
}
