package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	cmd := "emacs"
	args := append(
		[]string{"-nsl", "--no-site-file", "--no-splash"},
		os.Args[1:]...,
	)

	s, err := os.Stdin.Stat()
	if err != nil {
		fatal("Could not stat stdin: %v", err)
	}

	if (s.Mode() & os.ModeCharDevice) == 0 {
		var b bytes.Buffer
		n, err := io.Copy(&b, os.Stdin)
		if err != nil {
			fatal("Could not copy stdin to temp file: %v", err)
		}

		if n > 0 {
			args = append(
				args,
				[]string{"--eval", fmt.Sprintf("(insert \"%s\")", b.String())}...,
			)
		}
	}

	ecCmd := exec.Command(cmd, args...)
	ecCmd.Stderr = os.Stderr
	ecCmd.Stdout = os.Stdout

	if err = ecCmd.Start(); err != nil {
		fatal("Could run emacs: %v", err)
	}
}
