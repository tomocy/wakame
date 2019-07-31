package main

import (
	"fmt"
	"os"

	"github.com/tomocy/wakame/cmd/wakame/client"
)

func main() {
	client := client.NewCLI(os.Args)
	contri, err := client.FetchContributor()
	if err != nil {
		fmt.Println(err)
		client.ShowUsage()
		return
	}

	client.ShowContributor(contri)
}
