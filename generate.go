// +build ignore

package main

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

var mainTmpl = `
package main

var feedTmpl = ` + "`{{.}}`"

func main() {
	tmpl, err := template.New("template").Parse(mainTmpl)
	if err != nil {
		log.Fatalf("could not parse template: %v", err)
	}
	body, err := ioutil.ReadFile("template.html")
	if err != nil {
		log.Fatalf("could not open template: %v", err)
	}
	w, err := os.Create("template.go")
	if err != nil {
		log.Fatalf("could not create generated template file: %v", err)
	}
	if err := tmpl.Execute(w, string(body)); err != nil {
		log.Fatalf("could not execute template: %v", err)
	}
}
