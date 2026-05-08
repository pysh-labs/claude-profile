//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "FAIL" {
		fmt.Fprintln(os.Stderr, "simulated failure")
		os.Exit(1)
	}
	logPath := os.Getenv("FAKE_CLAUDE_LOG")
	if logPath == "" {
		os.Exit(0)
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Exit(2)
	}
	defer f.Close()
	fmt.Fprintln(f, strings.Join(os.Args[1:], " "))
}
