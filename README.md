# goread

Static RSS reader generator.

```bash
$ go get github.com/bake/goread
$ cd $GOPATH/src/github.com/bake/goread
$ make
$ echo "https://example.com/feed.rss" >> feeds.txt
$ echo "https://test.com/feed.xml" >> feeds.txt
$ goread # -in feeds.txt -out feeds.html
$
```
