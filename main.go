package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	s, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not stat stdin: %v", err)
		os.Exit(1)
	}

	cmd := "emacs"
	args := []string{"-nsl", "--no-site-file", "--no-splash"}

	if (s.Mode() & os.ModeCharDevice) == 0 {
		f, err := ioutil.TempFile("", "ec-")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create temp file: %v", err)
			os.Exit(1)
		}
		defer func() {
			f.Close()
			// os.Remove(f.Name())
		}()

		_, err = io.Copy(f, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not copy stdin to temp file: %v", err)
			os.Exit(1)
		}

		args = append(args, f.Name())
	}

	args = append(args, os.Args[1:]...)
	ecCmd := exec.Command(cmd, args...)
	ecCmd.Stderr = os.Stderr
	ecCmd.Stdout = os.Stdout

	if err = ecCmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Could run emacs: %v", err)
		os.Exit(1)
	}
}
