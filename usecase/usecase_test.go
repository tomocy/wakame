package usecase

import (
	"fmt"
	"testing"

	"github.com/tomocy/wakame/domain/model"
)

func TestFetch(t *testing.T) {
	mock := newMock()
	uc := NewContributorUsecase(mock)
	tests := []struct {
		name   string
		tester func(t *testing.T)
	}{
		{
			"success",
			func(t *testing.T) {
				repo := &model.Repository{Owner: "mock", Name: "mock"}
				expected := &model.Contributor{
					Name:          "alice",
					Repo:          repo,
					Contributions: 100,
				}
				contri, err := uc.Fetch(repo, "alice")
				if err != nil {
					t.Errorf("unexpected error returned by Fetch: %s\n", err)
				}
				if err := assertContributor(contri, expected); err != nil {
					t.Errorf("unexpected contributor returned by Fetch: %s\n", err)
				}
			},
		},
		{
			"no such contributor",
			func(t *testing.T) {
				repo := &model.Repository{Owner: "mock", Name: "mock"}
				contri, err := uc.Fetch(repo, "bob")
				if err == nil {
					t.Errorf("unexpected error returned by Fetch: got %s, want nil\n", err)
				}
				if contri != nil {
					t.Errorf("unexpected contributor returned by Fetch: got %v, want nil", contri)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.tester)
	}
}

func newMock() *mock {
	return &mock{
		repos: map[string][]*model.Contributor{
			"mock/mock": []*model.Contributor{
				{
					Name: "alice",
					Repo: &model.Repository{
						Owner: "mock",
						Name:  "mock",
					},
					Contributions: 100,
				},
			},
		},
	}
}

type mock struct {
	repos map[string][]*model.Contributor
}

func (m *mock) FetchContributors(repo *model.Repository, page int) ([]*model.Contributor, error) {
	name := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
	contris := m.repos[name]
	min := 10 * (page - 1)
	max := min + 10
	if len(contris) <= min {
		min = len(contris)
	}
	if len(contris) <= max {
		max = len(contris)
	}

	return contris[min:max], nil
}

func assertContributor(actual, expected *model.Contributor) error {
	if actual.Name != expected.Name {
		return fmt.Errorf("unexpected name of contributor: got %s, want %s", actual.Name, expected.Name)
	}
	if err := assertRepository(actual.Repo, expected.Repo); err != nil {
		return fmt.Errorf("unexpected repo of contributor: %s", err)
	}
	if actual.Contributions != expected.Contributions {
		return fmt.Errorf("unexpected contributions of contributor: got %d, want %d", actual.Contributions, expected.Contributions)
	}

	return nil
}

func assertRepository(actual, expected *model.Repository) error {
	if actual.Owner != expected.Owner {
		return fmt.Errorf("unexpected owner of repository: got %s, want %s", actual.Owner, expected.Owner)
	}
	if actual.Name != expected.Name {
		return fmt.Errorf("unexpected name of repository: got %s, want %s", actual.Name, expected.Name)
	}

	return nil
}
