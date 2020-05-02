// Package feed extends the package gofeed by adding a Date field to gofeed.Item
// which contains either the published or the updated date, since some feeds
// only offer a update time, as well as a sorting implementation based on the
// field.
package feed

import (
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

type Feed struct {
	gofeed.Feed
	Items []*Item
}

type Item struct {
	*gofeed.Item
	Feed     *gofeed.Feed
	Category string
}

func (i *Item) Time() time.Time {
	if i.PublishedParsed != nil {
		return *i.PublishedParsed
	}
	if i.UpdatedParsed != nil {
		return *i.UpdatedParsed
	}
	return time.Time{}
}

type SortByDate []*Item

func (is SortByDate) Len() int           { return len(is) }
func (is SortByDate) Less(i, j int) bool { return is[i].Time().Before(is[j].Time()) }
func (is SortByDate) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

type Parser struct{ gofeed.Parser }

func NewParser(c *http.Client) *Parser {
	p := gofeed.NewParser()
	p.Client = c
	return &Parser{*p}
}

func (p *Parser) ParseURL(url string) (*Feed, error) {
	f, err := p.Parser.ParseURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse feed")
	}
	items := make([]*Item, len(f.Items))
	for i, item := range f.Items {
		items[i] = &Item{Item: item, Feed: f}
	}
	return &Feed{*f, items}, err
}
