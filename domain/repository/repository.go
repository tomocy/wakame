package repository

import "github.com/tomocy/wakame/domain/model"

type ContributorRepository interface {
	FetchContributors(owner, repo string, page int) ([]*model.Contributor, error)
}
