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

func main() {
	feedPath := flag.String("feeds", "feeds.txt", "Path to a list of feed URLs")
	tmplPath := flag.String("template", "template.html", "Path to the HTML template")
	maxItems := flag.Int("max-items", 100, "Max number of items")
	flag.Parse()

	logger := log.New(os.Stderr, "", log.Lshortfile)
	client := http.DefaultClient

	body, err := ioutil.ReadFile(*feedPath)
	if err != nil {
		logger.Printf("could not read feeds from %s: %v", *feedPath, err)
	}
	urls := strings.Split(string(body), "\n")

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

	tmpl, err := template.New(path.Base(*tmplPath)).
		Funcs(template.FuncMap{
			"sanitize": bluemonday.StrictPolicy().Sanitize,
			"trim":     strings.TrimSpace,
		}).
		ParseFiles(*tmplPath)
	if err != nil {
		logger.Printf("could not parse template: %v", err)
	}

	data := struct {
		Items   []Item
		Updated time.Time
	}{items, time.Now()}
	if err := tmpl.Execute(os.Stdout, data); err != nil {
		logger.Printf("could not execute template: %v", err)
	}
}
