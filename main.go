package main

import (
	"fmt"
	"os"
	"github.com/jaymell/go-serve/daemon"
)

func run() {
	var d daemon.Daemon
	d.Init()
	d.Start()
}

func main() {
	if err := run(); err != nil {
		os.Stderr.WriteString(err)
		os.Exit(1)
	}
}
