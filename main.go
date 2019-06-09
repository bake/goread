//go:generate go run generate.go

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
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var version = "development"

func main() {
	inPath := flag.String("in", "feeds.yml", "Path to a list of feed URLs")
	outPath := flag.String("out", ".", "Path to generated HTML")
	tmplPath := flag.String("template", "", "Path to the HTML template")
	maxItems := flag.Int("max-items", 100, "Max number of items")
	concurrent := flag.Int64("n", 5, "Number of concurrent downloads")
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
			if len(str) <= 256 {
				return str
			}
			return str[:256] + " …"
		},
	}).Parse(feedTmpl))
	if *tmplPath != "" {
		tmpl, err = tmpl.ParseFiles(*tmplPath)
	}
	if err != nil {
		log.Fatalf("could not parse template: %v", err)
	}

	var allFeeds []*gofeed.Feed
	for cat, urls := range cats {
		var feeds []*gofeed.Feed
		feedc, errc := fetchAll(urls, *concurrent)
		for i := 0; i < len(urls); i++ {
			select {
			case feed := <-feedc:
				feeds = append(feeds, feed)
			case err := <-errc:
				log.Printf("could not fetch feed from %s: %v\n", cat, err)
			}
		}
		if len(feeds) > *maxItems {
			feeds = feeds[:*maxItems]
		}
		allFeeds = append(allFeeds, feeds...)
		if err := render(cat, catNames, feeds, tmpl, *outPath); err != nil {
			log.Printf("could not render %s: %v", cat, err)
		}
	}

	if len(allFeeds) > *maxItems {
		allFeeds = allFeeds[:*maxItems]
	}
	if err := render("index", catNames, allFeeds, tmpl, *outPath); err != nil {
		log.Printf("could not render index: %v", err)
	}
}

func render(category string, categories []string, feeds []*gofeed.Feed, tmpl *template.Template, outPath string) error {
	var items []item
	for _, f := range feeds {
		for _, i := range f.Items {
			t := time.Time{}
			if i.PublishedParsed != nil {
				t = *i.PublishedParsed
			}
			if i.UpdatedParsed != nil {
				t = *i.UpdatedParsed
			}
			items = append(items, item{*i, *f, t})
		}
	}
	sort.Sort(sort.Reverse(sortByPublished(items)))

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
