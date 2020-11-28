package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/subosito/gotenv"
)

const errProcessFinished = "os: process already finished"

func main() {
	path := flag.String("f", "", "filepath to .env")
	flag.Parse()

	env := parseEnv(*path)
	for key, value := range env {
		os.Setenv(key, value)
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Println("nothing to start")

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
		log.Fatal(err)
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
		case err = <-waitCh:
			if err == nil {
				os.Exit(0)
			}

			exitErr, ok := err.(*exec.ExitError)
			if ok {
				childStatus := exitErr.Sys().(syscall.WaitStatus)

				os.Exit(childStatus.ExitStatus())
			} else {
				log.Fatal(err)
			}
		case s := <-signCh:
			err = cmd.Process.Signal(s)
			if err != nil && err.Error() != errProcessFinished {
				log.Println("error sendind signal", s, err)
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
		log.Fatal(err)
	}

	// beacuse of we only read a file
	defer f.Close() //nolint:gosec

	env, err := gotenv.StrictParse(f)
	if err != nil {
		log.Fatal(err)
	}

	return env
}
