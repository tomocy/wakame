package client

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tomocy/wakame/domain/model"
	"github.com/tomocy/wakame/infra"
	"github.com/tomocy/wakame/usecase"
)

func NewCLI(args []string) *CLI {
	return &CLI{
		args: args,
	}
}

type CLI struct {
	args []string
}

func (c *CLI) Run() int {
	contri, err := c.FetchContributor()
	if err != nil {
		fmt.Println(err)
		c.ShowUsage()
		return 1
	}

	c.ShowContributor(contri)

	return 0
}

func (c *CLI) FetchContributor() (*model.Contributor, error) {
	config, err := c.parse()
	if err != nil {
		return nil, report("fetch contributor", err)
	}
	repo := new(infra.GitHub)
	uc := usecase.NewContributorUsecase(repo)

	return uc.Fetch(&model.Repository{
		Owner: config.owner,
		Name:  config.repo,
	}, config.uname)
}

func (c *CLI) parse() (*config, error) {
	config, err := c.parseConfig()
	if err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *CLI) parseConfig() (*config, error) {
	r := flag.String("r", "", "name of owner/repository")
	flag.Parse()
	splited := strings.Split(*r, "/")
	if len(splited) != 2 {
		return nil, errors.New("invalid format of name of owner/repository")
	}

	return &config{
		owner: splited[0],
		repo:  splited[1],
		uname: flag.Arg(0),
	}, nil
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
