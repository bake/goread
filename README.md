# goread

[![Go Report Card](https://goreportcard.com/badge/github.com/bake/goread)](https://goreportcard.com/report/github.com/bake/goread)

goread generates static HTML files showing previews of subscribed RSS feeds.
Similar to [rawdog](https://offog.org/code/rawdog/) and
[curn](http://software.clapper.org/curn/) it can be used as a cronjob.
Configuration happens through a simple YAML file containing categories and their
subscriptions. It comes with a built in template that can be overwritten.

See the [Screenshot](/screenshot.png).

## Basic Usage

Create a YAML file with the following format:

```yaml
golang:
  - https://blog.golang.org/feed.atom
  - https://campoy.cat/index.xml
  - https://medium.com/feed/@matryer

podcasts:
  - https://feeds.feedburner.com/SchrottcastTitusJonas
  - https://freakshow.fm/feed/m4a
```

By default, goread will look for a `feeds.yml` in the current directory and
renders its HTML files there too. This can be changed by using the `-in` and
`-out` flags. Note that `-in` expects a filename and `-out` a directory path.

```bash
$ goread -in ~/.goread.yml -out /var/www/html
$
```

## Use a custom template

The `-template` flag can be used to replace the defaut template. This repository
contains an additional one, [plain.html](/plain.html), which does not use any
external styling.

```bash
$ goread -template plain.html
$
```

## Help

```bash
$ goread -help
Usage of goread:
  -in string
        Path to a list of feed URLs (default "feeds.yml")
  -max-items int
        Max number of items (default 100)
  -n int
        Number of concurrent downloads (default 5)
  -out string
        Path to generated HTML (default ".")
  -template string
        Path to the HTML template
```

## Development

Use `go generate` to embed the default template. During development, you can use
the `-template` flag instead.
