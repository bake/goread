package main

import (
	"context"
	"net/http"

	"github.com/bake/goread/feed"
	"golang.org/x/sync/semaphore"
)

type request struct {
	cat, url string
}

type response struct {
	request
	feed *feed.Feed
	err  error
}

// fetch accepts the number of parallel downloads and returns a request and a
// response channel. The caller is responsible to close the request channel
// after all requests are enqueued, the response chan gets closed automatically.
func fetch(n int64, c *http.Client) (chan<- request, <-chan response) {
	sem := semaphore.NewWeighted(n)
	ctx := context.Background()
	reqc := make(chan request)
	resc := make(chan response)
	go func() {
		defer close(resc)
		defer sem.Acquire(ctx, n)
		for req := range reqc {
			sem.Acquire(ctx, 1)
			go func(req request) {
				defer sem.Release(1)
				feed, err := feed.NewParser(c).ParseURL(req.url)
				resc <- response{req, feed, err}
			}(req)
		}
	}()
	return reqc, resc
}

func fetchAll(n int64, fs feeds) <-chan response {
	reqc, resc := fetch(n, &http.Client{})
	go func() {
		defer close(reqc)
		for cat, urls := range fs {
			for _, url := range urls {
				reqc <- request{cat, url}
			}
		}
	}()
	return resc
}
