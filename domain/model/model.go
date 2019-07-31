package model

type Contributors []*Contributor

type Contributor struct {
	Name          string
	Contributions int
}
