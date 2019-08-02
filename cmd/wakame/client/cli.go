package client

import (
	"flag"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/tomocy/wakame/domain/model"
	"github.com/tomocy/wakame/infra"
	"github.com/tomocy/wakame/usecase"
)

func NewCLI() *CLI {
	return new(CLI)
}

type CLI struct{}

func (c *CLI) Run(config *Config) int {
	contri, err := c.FetchContributor(config)
	if err != nil {
		fmt.Println(err)
		c.ShowUsage()
		return 1
	}

	c.ShowContributor(contri)

	return 0
}

func (c *CLI) FetchContributor(config *Config) (*model.Contributor, error) {
	repo := new(infra.GitHub)
	uc := usecase.NewContributorUsecase(repo)

	return uc.Fetch(&model.Repository{
		Owner: config.Owner,
		Name:  config.Repo,
	}, config.Username)
}

func (c *CLI) ShowUsage() {
	flag.Usage()
}

func (c *CLI) ShowContributor(contri *model.Contributor) {
	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"Username", "Repository", "Contributions"})
	w.Append([]string{contri.Name, contri.Repo.String(), fmt.Sprintf("%d", contri.Contributions)})

	w.Render()
}
