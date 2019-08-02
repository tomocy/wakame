package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tomocy/wakame/cmd/wakame/client"
	"github.com/urfave/cli"
)

func main() {
	app := newApp()
	if err := app.run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run app: %s\n", err)
		return
	}
}

func newApp() *app {
	a := new(app)
	a.setUp()

	return a
}

type app struct {
	driver *cli.App
}

func (a *app) setUp() {
	a.driver = cli.NewApp()
	a.setBasic()
	a.setCommands()
}

func (a *app) setBasic() {
	a.driver.Name = "wakame"
}

func (a *app) setCommands() {
	a.driver.Commands = []cli.Command{
		{
			Name:   "cli",
			Action: a.runCLI,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "r",
				},
			},
		},
		{
			Name:   "html",
			Action: a.runHTML,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "a",
					Value: ":80",
				},
			},
		},
	}
}

func (a *app) runCLI(ctx *cli.Context) error {
	c := client.NewCLI()
	splited := strings.SplitN(ctx.String("r"), "/", 2)
	config := &client.Config{
		Owner:    splited[0],
		Repo:     splited[1],
		Username: ctx.Args().First(),
	}
	if err := config.Validate(); err != nil {
		return err
	}

	return c.Run(config)
}

func (a *app) runHTML(ctx *cli.Context) error {
	c := client.NewHTML(ctx.String("a"))

	return c.Run()
}

func (a *app) run(args []string) error {
	return a.driver.Run(args)
}
