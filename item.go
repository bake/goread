package main

import "github.com/mmcdole/gofeed"

type Item struct {
	gofeed.Item
	Feed gofeed.Feed
}

type SortByPublished []Item

func (is SortByPublished) Len() int { return len(is) }
func (is SortByPublished) Less(i, j int) bool {
	return is[i].PublishedParsed.Before(*is[j].PublishedParsed)
}
func (is SortByPublished) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
