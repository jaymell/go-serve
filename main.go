package main

import (
	"fmt"
	"github.com/jaymell/go-serve/daemon"
	"github.com/jaymell/go-serve/api"
	"os"
)

func run() error {
    var d daemon.Daemon

    f, err := os.Open("config.api.json")
    if err != nil {
        return fmt.Errorf("unable to open api config file: ", err)
    }

    api, err := api.New(f)
    if err != nil {
        return fmt.Errorf("failed to initialize api: ", err)
    }

    f, err = os.Open("config.daemon.json")
    if err != nil {
        return fmt.Errorf("unable to open daemon config file: ", err)
    }

    err = d.Init(api, f)
    if err != nil {
        return fmt.Errorf("failed to initialize daemon: ", err)

    }

    d.Start()

    return nil
}

func main() {
    fmt.Println("Starting server... ")
    err := run()
    if err != nil {
        fmt.Println("failed: ", err)
        os.Exit(1)
    }
}

