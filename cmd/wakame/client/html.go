package client

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/tomocy/wakame/domain/model"
	"github.com/tomocy/wakame/infra"
	"github.com/tomocy/wakame/usecase"

	"github.com/go-chi/chi"
	"github.com/tomocy/caster"
)

func NewHTML(addr string) *HTML {
	return &HTML{
		addr: addr,
	}
}

type HTML struct {
	addr   string
	caster caster.Caster
}

func (h *HTML) Run() error {
	if err := h.load(); err != nil {
		return err
	}

	hand := h.prepareHandler()
	if err := http.ListenAndServe(h.addr, hand); err != nil {
		return err
	}

	return nil
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
		Owner: config.Owner,
		Name:  config.Repo,
	}, config.Username)
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

func (h *HTML) parse(r *http.Request) (*Config, error) {
	config := &Config{
		Owner:    chi.URLParam(r, "owner"),
		Repo:     chi.URLParam(r, "repo"),
		Username: chi.URLParam(r, "uname"),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
