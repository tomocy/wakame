package main

import (
	"os"

	"github.com/tomocy/wakame/cmd/wakame/client"
)

func main() {
	runner := newRunner()
	os.Exit(runner.Run())
}

type app struct{}

func newRunner() runner {
	return client.NewCLI(os.Args)
}

type runner interface {
	Run() int
}
