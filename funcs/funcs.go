// Package funcs contains a collection of generic templating functions.
package funcs

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

func FuncMap(truncateLen int) template.FuncMap {
	return template.FuncMap{
		"hash":     Hash(),
		"sanitize": Sanitize(),
		"time":     Time(),
		"title":    Title(),
		"trim":     Trim(),
		"truncate": Truncate(truncateLen),
	}
}

func Hash() func(string) string {
	return func(s string) string { return fmt.Sprintf("%x", sha1.Sum([]byte(s))) }
}

func Sanitize() func(string) string { return bluemonday.StrictPolicy().Sanitize }

func Time() func() time.Time { return func() time.Time { return time.Now() } }

func Title() func(string) string { return strings.Title }

func Trim() func(string) string { return strings.TrimSpace }

func Truncate(n int) func(string) string {
	ellipsis := " …"
	return func(s string) string {
		if len(s)-len(ellipsis) <= n {
			return s
		}
		return s[:n] + ellipsis
	}
}
