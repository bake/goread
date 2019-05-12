//go:generate go run generate.go

package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

var version = "development"

func main() {
	feedPath := flag.String("feeds", "feeds.txt", "Path to a list of feed URLs")
	tmplPath := flag.String("template", "", "Path to the HTML template")
	maxItems := flag.Int("max-items", 100, "Max number of items")
	flag.Parse()

	logger := log.New(os.Stderr, "", log.Lshortfile)
	client := http.DefaultClient

	body, err := ioutil.ReadFile(*feedPath)
	if err != nil {
		logger.Fatalf("could not read feeds from %s: %v", *feedPath, err)
	}
	var urls []string
	for _, url := range strings.Split(string(body), "\n") {
		if url != "" {
			urls = append(urls, url)
		}
	}

	feeds, err := fetchAll(client, urls)
	if err != nil {
		logger.Printf("could not fetch feeds: %v", err)
	}
	var items []Item
	for _, f := range feeds {
		for _, i := range f.Items {
			items = append(items, Item{*i, *f})
		}
	}
	sort.Sort(sort.Reverse(SortByPublished(items)))
	if len(items) > *maxItems {
		items = items[:*maxItems]
	}

	tmpl, err := template.New(path.Base(*tmplPath)).Funcs(template.FuncMap{
		"sanitize": bluemonday.StrictPolicy().Sanitize,
		"trim":     strings.TrimSpace,
	}).Parse(feedTmpl)
	if *tmplPath != "" {
		tmpl, err = tmpl.ParseFiles(*tmplPath)
	}
	if err != nil {
		logger.Fatalf("could not parse template: %v", err)
	}

	data := struct {
		Items   []Item
		Updated time.Time
		Version string
	}{items, time.Now(), version}
	if err := tmpl.Execute(os.Stdout, data); err != nil {
		logger.Fatalf("could not execute template: %v", err)
	}
}
