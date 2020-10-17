package main

import (
	"bytes"
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
		var buf bytes.Buffer
		_, err := io.Copy(&buf, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not copy stdin to temp file: %v", err)
			os.Exit(1)
		}

		f, err := ioutil.TempFile("", "ec-")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create temp file: %v", err)
			os.Exit(1)
		}
		defer func() {
			f.Close()
			// os.Remove(f.Name())
		}()

		buf.WriteTo(f)
		args = append(args, f.Name())
	}

	args = append(args, os.Args[1:]...)
	ecCmd := exec.Command(cmd, args...)
	err = run(ecCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could run emacs: %v", err)
		os.Exit(1)
	}
}

func run(cmd *exec.Cmd) error {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Start()
}
