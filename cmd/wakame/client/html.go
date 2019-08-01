package client

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/tomocy/caster"
)

type HTML struct {
	caster caster.Caster
}

func (h *HTML) load() error {
	var err error
	h.caster, err = caster.New(&caster.TemplateSet{
		Filenames: []string{h.join("master.html")},
	})

	return err
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

func (h *HTML) register(r chi.Router) {}
