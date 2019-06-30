package main

import (
	"html/template"
	"os"
	"path"
	"time"

	"github.com/bake/goread/feed"
	"github.com/pkg/errors"
)

type page struct {
	tmpl *template.Template
	out  string
	max  int

	Category   string
	Categories []string
	Items      []*feed.Item
	Updated    time.Time
	Version    string
}

func (p *page) render(name, category string, items []*feed.Item) error {
	p.Items = items
	p.Category = category
	p.Updated = time.Now()
	w, err := os.Create(path.Join(p.out, name+".html"))
	if len(p.Items) > p.max {
		p.Items = p.Items[:p.max]
	}
	if err != nil {
		return errors.Wrap(err, "could not generate output file")
	}
	if err := p.tmpl.Execute(w, p); err != nil {
		return errors.Wrap(err, "could not execute template")
	}
	return nil
}
