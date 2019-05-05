# goread

Static RSS reader generator.

```bash
$ go get git.192k.pw/bake/goread
$ goread -help
Usage of goread:
  -feeds string
    Path to a list of feed URLs (default "feeds.txt")
  -max-items int
    Max number of items (default 100)
  -template string
    Path to the HTML template (default "template.html")
$ wget https://raw.githubusercontent.com/bake/goread/master/template.html
$ echo "https://example.com/feed.rss" >> feeds.txt
$ echo "https://test.com/feed.xml" >> feeds.txt
$ goread > feed.html 2>> goread.log
$
```
