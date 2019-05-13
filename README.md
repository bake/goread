# goread

[![Go Report Card](https://goreportcard.com/badge/github.com/bake/goread)](https://goreportcard.com/report/github.com/bake/goread)

goread generates static HTML files showing previews of subscribed RSS feeds.
Subscriptions are handled through a textfile containing one feed URL per line.
The default template can be overwritten with a flag.

It can be used as a cronjob.

```bash
$ go get github.com/bake/goread
$ cd $GOPATH/src/github.com/bake/goread
$ make
$
```

```bash
$ echo "https://example.com/feed.rss" > feeds.txt
$ echo "https://test.com/feed.xml" >> feeds.txt
$ ./goread # -in feeds.txt -out feeds.html
$
```
