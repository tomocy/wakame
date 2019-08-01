package infra

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tomocy/wakame/domain/model"
)

type GitHub struct{}

func (gh *GitHub) FetchContributors(repo *model.Repository, page int) ([]*model.Contributor, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors?page=%d", repo.Owner, repo.Name, page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fetcheds contributors
	if err := json.NewDecoder(resp.Body).Decode(&fetcheds); err != nil {
		return nil, err
	}
	adapteds := fetcheds.adapt()
	for i := range adapteds {
		adapteds[i].Repo = repo
	}

	return adapteds, nil
}

type contributors []*contributor

func (cs contributors) adapt() []*model.Contributor {
	adapteds := make([]*model.Contributor, len(cs))
	for i, c := range cs {
		adapteds[i] = c.adapt()
	}

	return adapteds
}

type contributor struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

func (c *contributor) adapt() *model.Contributor {
	return &model.Contributor{
		Name:          c.Login,
		Contributions: c.Contributions,
	}
}
