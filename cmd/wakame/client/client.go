package client

import (
	"errors"
	"fmt"
)

type Config struct {
	Owner, Repo string
	Username    string
}

func (c *Config) Validate() error {
	if c.Owner == "" {
		return errors.New("owner is empty")
	}
	if c.Repo == "" {
		return errors.New("repo is empty")
	}
	if c.Username == "" {
		return errors.New("username is empty")
	}

	return nil
}

func report(did string, err error) error {
	return fmt.Errorf("failed to %s: %s", did, err)
}
