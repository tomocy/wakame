package repository

import "github.com/tomocy/wakame/domain/model"

type ContributorRepository interface {
	FetchContributors(repo *model.Repository, page int) ([]*model.Contributor, error)
}
