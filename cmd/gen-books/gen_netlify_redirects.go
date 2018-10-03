package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	// https://www.netlify.com/docs/headers-and-basic-auth/#custom-headers
	netlifyHeaders = `
# long-lived caching
/s/*
  Cache-Control: max-age=31536000
/*
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  X-XSS-Protection: 1; mode=block
`
)

func genNetlifyHeaders() {
	path := filepath.Join("www", "_headers")
	err := ioutil.WriteFile(path, []byte(netlifyHeaders), 0644)
	panicIfErr(err)
}

func genNetlifyRedirectsForBook(b *Book) []string {
	var res []string

	// TODO: this should be recursive. Maybe add a way to iterate over all pages
	// in a book
	for _, page := range b.RootPage.Pages {
		id := page.NotionID
		s := fmt.Sprintf(`/essential/%s/%s* /essential/%s/404.html 404`, b.Dir, id, b.Dir)
		// TODO: also add redirects for old article ids (only needed for Go book)
		// alternatively just forget about it
		res = append(res, s)
	}

	// catch-all redirect for all other missing pages
	s := fmt.Sprintf(`/essential/%s/* /essential/%s/404.html 404`, b.Dir, b.Dir)
	res = append(res, s)
	res = append(res, "")
	return res
}

func genNetlifyRedirects() {
	var a []string
	for _, b := range books {
		ab := genNetlifyRedirectsForBook(b)
		a = append(a, ab...)
	}
	s := strings.Join(a, "\n")
	path := filepath.Join("www", "_redirects")
	err := ioutil.WriteFile(path, []byte(s), 0644)
	panicIfErr(err)
}
