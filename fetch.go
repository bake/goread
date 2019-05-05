package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func fetch(c *http.Client, url string) (*gofeed.Feed, error) {
	res, err := c.Get(url)
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

func fetchAll(c *http.Client, urls []string) ([]*gofeed.Feed, error) {
	con := int64(2)
	sem := semaphore.NewWeighted(int64(con))
	ctx := context.Background()
	var g errgroup.Group
	feedc := make(chan *gofeed.Feed)
	go func() {
		defer close(feedc)
		for _, url := range urls {
			sem.Acquire(ctx, 1)
			url := url
			g.Go(func() error {
				defer sem.Release(1)
				feed, err := fetch(c, url)
				if err != nil {
					return errors.Wrapf(err, "could not get %s", url)
				}
				feedc <- feed
				return nil
			})
		}
		sem.Acquire(ctx, con)
	}()
	var feeds []*gofeed.Feed
	for feed := range feedc {
		feeds = append(feeds, feed)
	}
	if err := g.Wait(); err != nil {
		log.Println(err)
	}
	return feeds, nil
}
