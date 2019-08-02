package client

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	fmt.Printf("listen and serve on %s\n", h.addr)
	if err := http.ListenAndServe(h.addr, hand); err != nil {
		return err
	}

	return nil
}

func (h *HTML) load() error {
	var err error
	h.caster, err = caster.New(&caster.TemplateSet{
		Filenames: []string{h.joinHTML("master.html")},
	})
	if err != nil {
		return err
	}

	if err := h.caster.ExtendAll(map[string]*caster.TemplateSet{
		"contributor.new": &caster.TemplateSet{
			Filenames: []string{h.joinHTML("contributor/new.html")},
		},
		"contributor.single": &caster.TemplateSet{
			Filenames: []string{h.joinHTML("contributor/single.html")},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (h *HTML) joinHTML(ss ...string) string {
	ps := append([]string{"html"}, ss...)
	return h.join(ps...)
}

func (h *HTML) join(ss ...string) string {
	dir := filepath.Join(os.Getenv("GOPATH"), "src/github.com/tomocy/wakame/cmd/wakame/client/resource")
	ps := append([]string{dir}, ss...)
	return filepath.Join(ps...)
}

func (h *HTML) prepareHandler() http.Handler {
	r := chi.NewRouter()
	h.register(r)

	return r
}

func (h *HTML) register(r chi.Router) {
	r.Get("/css/*", http.StripPrefix("/css", http.FileServer(http.Dir(h.join("css")))).ServeHTTP)
	r.Get("/", h.showContributorOrSearchForm)
}

func (h *HTML) showContributorOrSearchForm(w http.ResponseWriter, r *http.Request) {
	config, err := h.parse(r)
	if err != nil {
		h.showContributorSearchForm(w, config)
		return
	}

	h.showContributor(w, config)
}

func (h *HTML) showContributorSearchForm(w http.ResponseWriter, config *Config) {
	if err := h.caster.Cast(w, "contributor.new", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *HTML) showContributor(w http.ResponseWriter, config *Config) {
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
	q := r.URL.Query()
	repo := "/"
	if q.Get("r") != "" {
		repo = q.Get("r")
	}
	splited := strings.SplitN(repo, "/", 2)
	config := &Config{
		Owner:    splited[0],
		Repo:     splited[1],
		Username: q.Get("u"),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
