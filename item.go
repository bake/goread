package main

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type item struct {
	gofeed.Item
	Feed gofeed.Feed
	Time time.Time
}

type sortByPublished []item

func (is sortByPublished) Len() int           { return len(is) }
func (is sortByPublished) Less(i, j int) bool { return is[i].Time.Before(is[j].Time) }
func (is sortByPublished) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }
