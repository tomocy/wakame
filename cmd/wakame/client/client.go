package client

import (
	"errors"
	"fmt"
)

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
