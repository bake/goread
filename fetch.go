package main

import (
	"context"
	"net/http"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

func fetch(url string) (*gofeed.Feed, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get feed at %s", url)
	}
	defer res.Body.Close()
	fp := gofeed.NewParser()
	f, err := fp.Parse(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse feed at %s", url)
	}
	return f, nil
}

func fetchAll(urls []string, n int64) (chan *gofeed.Feed, chan error) {
	sem := semaphore.NewWeighted(n)
	ctx := context.Background()
	feedc := make(chan *gofeed.Feed)
	errc := make(chan error)
	go func() {
		defer close(errc)
		defer close(feedc)
		for _, url := range urls {
			sem.Acquire(ctx, 1)
			url := url
			go func() {
				feed, err := fetch(url)
				if err != nil {
					errc <- err
					return
				}
				feedc <- feed
			}()
		}
		sem.Acquire(ctx, n)
	}()
	return feedc, errc
}
