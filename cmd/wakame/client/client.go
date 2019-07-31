package client

import (
	"errors"
	"flag"
	"fmt"

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

	return uc.Fetch(config.owner, config.repo, config.uname)
}

func (c *CLI) parse() (*config, error) {
	parsed := new(config)
	flag.StringVar(&parsed.owner, "owner", "", "name of owner")
	flag.StringVar(&parsed.repo, "repo", "", "name of repository")
	flag.Parse()
	parsed.uname = flag.Arg(0)

	if err := parsed.validate(); err != nil {
		return nil, err
	}

	return parsed, nil
}

func (c *CLI) ShowUsage() {
	flag.Usage()
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
