package client

import (
	"os"
	"path/filepath"

	"github.com/tomocy/caster"
)

func NewHTML() *HTML {
	return new(HTML)
}

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
