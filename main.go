//go:generate go run generate.go
//go:generate go fmt template.go

package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"path"
	"sort"

	"github.com/bake/goread/feed"
	"github.com/bake/goread/funcs"
	"gopkg.in/yaml.v2"
)

var version = "development"

type feeds map[string][]string

func main() {
	inPath := flag.String("in", "feeds.yml", "Path to a list of feed URLs")
	outPath := flag.String("out", ".", "Path to generated HTML")
	tmplPath := flag.String("template", "", "Path to the HTML template")
	maxItems := flag.Int("max-items", 100, "Max number of items per page")
	concurrent := flag.Int64("n", 5, "Number of concurrent downloads")
	truncateLen := flag.Int("truncate-length", 256, "Number of characters per feed item")
	flag.Parse()

	p := page{
		out:     *outPath,
		max:     *maxItems,
		Version: version,
	}

	var err error
	p.tmpl, err = template.
		New(path.Base(*tmplPath)).
		Funcs(funcs.FuncMap(*truncateLen)).
		Parse(feedTmpl)
	if err != nil {
		log.Fatalf("could not parse internal template: %v", err)
	}
	if *tmplPath != "" {
		p.tmpl, err = p.tmpl.ParseFiles(*tmplPath)
	}
	if err != nil {
		log.Fatalf("could not parse template: %v", err)
	}

	r, err := os.Open(*inPath)
	if err != nil {
		log.Fatalf("could not open feeds: %v", err)
	}
	defer r.Close()
	var fs feeds
	if err := yaml.NewDecoder(r).Decode(&fs); err != nil {
		log.Fatalf("could not decode %s: %v", path.Base(*inPath), err)
	}

	var items []*feed.Item
	for res := range fetchAll(*concurrent, fs) {
		if res.err != nil {
			log.Printf("could not get %s: %v", res.url, res.err)
			continue
		}
		for _, item := range res.feed.Items {
			item.Category = res.cat
			items = append(items, item)
		}
	}
	sort.Sort(sort.Reverse(feed.SortByDate(items)))

	cats := map[string][]*feed.Item{"index": items}
	feeds := map[string][]*feed.Item{}
	hash := funcs.Hash()
	for _, item := range items {
		cats[item.Category] = append(cats[item.Category], item)
		feeds[hash(item.Feed.Link)] = append(feeds[hash(item.Feed.Link)], item)
	}
	for cat := range cats {
		p.Categories = append(p.Categories, cat)
	}
	sort.Strings(p.Categories)
	for cat, items := range cats {
		if err := p.render(cat, cat, items); err != nil {
			log.Fatalf("could not render %s: %v", cat, err)
		}
	}
	for feed, items := range feeds {
		if err := p.render(feed, items[0].Feed.Title, items); err != nil {
			log.Fatalf("could not render %s: %v", feed, err)
		}
	}
}
