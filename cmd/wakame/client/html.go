package client

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tomocy/wakame/domain/model"
	"github.com/tomocy/wakame/infra"
	"github.com/tomocy/wakame/usecase"

	"github.com/go-chi/chi"
	"github.com/tomocy/caster"
)

func NewHTML() *HTML {
	return new(HTML)
}

type HTML struct {
	caster caster.Caster
}

func (h *HTML) Run() int {
	if err := h.load(); err != nil {
		fmt.Printf("failed for html to run: %s\n", err)
		return 1
	}

	hand := h.prepareHandler()
	if err := http.ListenAndServe(":80", hand); err != nil {
		fmt.Printf("failed for html to run: %s\n", err)
		return 1
	}

	return 0
}

func (h *HTML) load() error {
	var err error
	h.caster, err = caster.New(&caster.TemplateSet{
		Filenames: []string{h.join("master.html")},
	})
	if err != nil {
		return err
	}

	if err := h.caster.ExtendAll(map[string]*caster.TemplateSet{
		"contributor.single": &caster.TemplateSet{
			Filenames: []string{h.join("contributor/single.html")},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (h *HTML) join(ss ...string) string {
	dir := filepath.Join(os.Getenv("GOPATH"), "src/github.com/tomocy/wakame/cmd/wakame/client/resource/html")
	ps := append([]string{dir}, ss...)
	return filepath.Join(ps...)
}

func (h *HTML) prepareHandler() http.Handler {
	r := chi.NewRouter()
	h.register(r)

	return r
}

func (h *HTML) register(r chi.Router) {
	r.Get("/{owner}/{repo}/{uname}", h.showContributor)
}

func (h *HTML) showContributor(w http.ResponseWriter, r *http.Request) {
	config, err := h.parse(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	repo := new(infra.GitHub)
	uc := usecase.NewContributorUsecase(repo)
	contri, err := uc.Fetch(&model.Repository{
		Owner: config.owner,
		Name:  config.repo,
	}, config.uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.caster.Cast(w, "contributor.single", map[string]interface{}{
		"Contributor": contri,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *HTML) parse(r *http.Request) (*config, error) {
	config := &config{
		owner: chi.URLParam(r, "owner"),
		repo:  chi.URLParam(r, "repo"),
		uname: chi.URLParam(r, "uname"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}
