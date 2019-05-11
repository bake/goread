# goread

Static RSS reader generator.

```bash
$ go get github.com/bake/goread
$ cd $GOPATH/src/github.com/bake/goread
$ go generate
$ go install
$ echo "https://example.com/feed.rss" >> feeds.txt
$ echo "https://test.com/feed.xml" >> feeds.txt
$ goread > feed.html 2>> goread.log
$
```
