package client

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/tomocy/wakame/domain/model"
	"github.com/tomocy/wakame/infra/github"
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

func (c *CLI) FetchContributor() (*model.Contributor, error) {
	config, err := c.parse()
	if err != nil {
		return nil, report("fetch contributor", err)
	}
	repo := new(github.GitHub)
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
	w.SetHeader([]string{"Username", "Contributions"})
	w.Append([]string{contri.Name, fmt.Sprintf("%d", contri.Contributions)})

	w.Render()
}

type config struct {
	owner, repo string
	uname       string
}

func (c *config) validate() error {
	if c.owner == "" {
		return errors.New("owner is empty")
	}
	if c.repo == "" {
		return errors.New("repo is empty")
	}
	if c.uname == "" {
		return errors.New("username is empty")
	}

	return nil
}

func report(did string, err error) error {
	return fmt.Errorf("failed to %s: %s", did, err)
}
