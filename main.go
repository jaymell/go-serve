package main

import (
	"fmt"
	"os"
	"github.com/jaymell/go-serve/daemon"
)

func run() error {
	var d daemon.Daemon
	err := d.Init()
	if err != nil {
		return fmt.Errorf("failed on daemon init: ", err)
	}
	d.Start()
	return nil
}

func main() {
	err := run()
	if err != nil {
	 	fmt.Println("failed: ", err)
	 	os.Exit(1)
	}
}
