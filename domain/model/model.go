package model

import "errors"

type Contributors []*Contributor

func (cs Contributors) Find(name string) (*Contributor, error) {
	for _, c := range cs {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, errors.New("no such contributor")
}

type Contributor struct {
	Name          string
	Repo          *Repository
	Contributions int
}

type Repository struct {
	Owner, Name string
}
