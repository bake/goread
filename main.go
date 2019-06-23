//go:generate go run generate.go
//go:generate go fmt template.go

package main

import (
	"flag"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var version = "development"

func main() {
	inPath := flag.String("in", "feeds.yml", "Path to a list of feed URLs")
	outPath := flag.String("out", ".", "Path to generated HTML")
	tmplPath := flag.String("template", "", "Path to the HTML template")
	maxItems := flag.Int("max-items", 100, "Max number of items per page")
	concurrent := flag.Int64("n", 5, "Number of concurrent downloads")
	truncateLen := flag.Int("truncate-length", 256, "Number of characters per feed item")
	flag.Parse()

	r, err := os.Open(*inPath)
	if err != nil {
		log.Fatalf("could not open feeds: %v", err)
	}
	var cats map[string][]string
	if err := yaml.NewDecoder(r).Decode(&cats); err != nil {
		log.Fatal(err)
	}

	var catNames []string
	for cat := range cats {
		catNames = append(catNames, cat)
	}
	sort.Strings(catNames)

	tmpl := template.Must(template.New(path.Base(*tmplPath)).Funcs(template.FuncMap{
		"sanitize": bluemonday.StrictPolicy().Sanitize,
		"trim":     strings.TrimSpace,
		"truncate": func(str string) string {
			if len(str) <= *truncateLen {
				return str
			}
			return str[:*truncateLen] + " …"
		},
	}).Parse(feedTmpl))
	if *tmplPath != "" {
		tmpl, err = tmpl.ParseFiles(*tmplPath)
	}
	if err != nil {
		log.Fatalf("could not parse template: %v", err)
	}

	var allItems []item
	for cat, urls := range cats {
		var items []item
		feedc, errc := fetchAll(urls, *concurrent)
		for range urls {
			select {
			case feed := <-feedc:
				for _, item := range feed.Items {
					items = append(items, newItem(item, feed))
				}
			case err := <-errc:
				log.Printf("could not fetch feed from %s: %v\n", cat, err)
			}
		}
		sort.Sort(sort.Reverse(sortByPublished(items)))
		if len(items) > *maxItems {
			items = items[:*maxItems]
		}
		allItems = append(allItems, items...)
		if err := render(cat, catNames, items, tmpl, *outPath); err != nil {
			log.Printf("could not render %s: %v", cat, err)
		}
	}

	sort.Sort(sort.Reverse(sortByPublished(allItems)))
	if len(allItems) > *maxItems {
		allItems = allItems[:*maxItems]
	}
	if err := render("index", catNames, allItems, tmpl, *outPath); err != nil {
		log.Printf("could not render index: %v", err)
	}
}

func render(category string, categories []string, items []item, tmpl *template.Template, outPath string) error {
	data := struct {
		Category   string
		Categories []string
		Items      []item
		Updated    time.Time
		Version    string
	}{category, categories, items, time.Now(), version}
	w, err := os.Create(path.Join(outPath, category+".html"))
	if err != nil {
		return errors.Wrap(err, "could not generate output file")
	}
	if err := tmpl.Execute(w, data); err != nil {
		return errors.Wrap(err, "could not execute template")
	}
	return nil
}
