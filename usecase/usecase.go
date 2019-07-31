package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/tomocy/wakame/domain/model"
	"github.com/tomocy/wakame/domain/repository"
)

func NewContributorUsecase(repo repository.ContributorRepository) *ContributorUsecase {
	return &ContributorUsecase{
		repo: repo,
	}
}

type ContributorUsecase struct {
	repo repository.ContributorRepository
}

func (u *ContributorUsecase) Fetch(owner, repo, uname string) (*model.Contributor, error) {
	for page := 1; ; page++ {
		fetcheds, err := u.repo.FetchContributors(owner, repo, page)
		if err != nil {
			return nil, report("fetch contributors", err)
		}
		if len(fetcheds) <= 0 {
			break
		}

		found, err := (model.Contributors(fetcheds)).Find(uname)
		if err == nil {
			return found, nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil, errors.New("no such contributor")
}

func report(did string, err error) error {
	return fmt.Errorf("failed to %s: %s", did, err)
}
