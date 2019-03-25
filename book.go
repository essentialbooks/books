package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/kjk/u"
)

// Book represents a book
type Book struct {
	Title     string // "Go", "jQuery" etcc
	TitleLong string // "Essential Go", "Essential jQuery" etc.

	NotionStartPageID string
	RootPage          *Page   // a tree of pages
	cachedPages       []*Page // pages flattened into an array

	idToPage map[string]*Page

	Dir            string // directory name for the book e.g. "go"
	SoContributors []SoContributor

	defaultLang string // default programming language for programming examples
	knownUrls   []string

	// generated toc javascript data
	tocData []byte
	// url of combined tocData and app.js
	AppJSURL string

	// name of a file in covers/ directory
	// e.g. "Python.png"
	CoverImageName string

	// cache related
	cache *Cache

	// for concurrency
	sem chan bool
	wg  sync.WaitGroup
}

// CacheDir returns a cache dir for this book
func (b *Book) CacheDir() string {
	panicIf(b.Dir == "", "b.Dir should not be empty")
	return filepath.Join("cache", b.Dir)
}

// OutputCacheDir returns output cache dir for this book
func (b *Book) OutputCacheDir() string {
	return filepath.Join(b.CacheDir(), "output")
}

// NotionCacheDir returns output cache dir for this book
func (b *Book) NotionCacheDir() string {
	return filepath.Join(b.CacheDir(), "notion")
}

// ReplitCachePath returns path of the cache file for replits
func (b *Book) ReplitCachePath() string {
	return filepath.Join(b.CacheDir(), "replit_cache.txt")
}

func (b *Book) CachePath() string {
	return filepath.Join(b.CacheDir(), "cache.txt")
}

// SourceDir is where source files for a given book are
func (b *Book) SourceDir() string {
	return filepath.Join("books", b.Dir)
}

// this is where html etc. files for a book end up
func (b *Book) destDir() string {
	return filepath.Join(destEssentialDir, b.Dir)
}

// ContributorCount returns number of contributors
func (b *Book) ContributorCount() int {
	return len(b.SoContributors)
}

// ContributorsURL returns url of the chapter that lists contributors
func (b *Book) ContributorsURL() string {
	return b.URL() + "/9999-contributors"
}

// GitHubURL returns link to GitHub for this book
func (b *Book) GitHubURL() string {
	return gitHubBaseURL + "/tree/master/books/" + filepath.Base(b.destDir())
}

// URL returns url of the book, used in index.tmpl.html
func (b *Book) URL() string {
	return fmt.Sprintf("/essential/%s/", b.Dir)
}

// CanonnicalURL returns full url including host
func (b *Book) CanonnicalURL() string {
	return urlJoin(siteBaseURL, b.URL())
}

// ShareOnTwitterText returns text for sharing on twitter
func (b *Book) ShareOnTwitterText() string {
	return fmt.Sprintf(`"%s" - a free programming book`, b.TitleLong)
}

// CoverURL returns url to cover image
func (b *Book) CoverURL() string {
	panicIf(b.CoverImageName == "")
	return fmt.Sprintf("/covers/%s", b.CoverImageName)
}

// CoverFullURL returns a URL for the cover including host
func (b *Book) CoverFullURL() string {
	return urlJoin(siteBaseURL, b.CoverURL())
}

// CoverTwitterFullURL returns a URL for the cover including host
func (b *Book) CoverTwitterFullURL() string {
	panicIf(b.CoverImageName == "")
	coverURL := fmt.Sprintf("/covers/twitter/%s", b.CoverImageName)
	return urlJoin(siteBaseURL, coverURL)
}

// Chapters returns pages that are top-level chapters
func (b *Book) Chapters() []*Page {
	return b.RootPage.Pages
}

// GetAllPages returns all pages, flattened
func (b *Book) GetAllPages() []*Page {
	// to prevent infinite recursion if pages show up twice (shouldn't happen)
	if len(b.cachedPages) > 0 {
		return b.cachedPages
	}
	if b.RootPage == nil {
		return nil
	}
	seen := map[*Page]bool{}
	pages := []*Page{b.RootPage}
	curr := 0
	for curr < len(pages) {
		page := pages[curr]
		curr++
		if seen[page] {
			continue
		}
		seen[page] = true
		for _, p := range page.Pages {
			p.Parent = page
		}
		pages = append(pages, page.Pages...)
	}
	return pages
}

// PagesCount returns total number of articles
func (b *Book) PagesCount() int {
	return len(b.GetAllPages()) - 1 // don't count top page
}

// ChaptersCount returns number of chapters
func (b *Book) ChaptersCount() int {
	return len(b.RootPage.Pages)
}

func updateBookAppJS(book *Book) {
	srcName := fmt.Sprintf("app-%s.js", book.Dir)
	path := filepath.Join("tmpl", "app.js")
	d, err := ioutil.ReadFile(path)
	maybePanicIfErr(err)
	if err != nil {
		return
	}
	if doMinify {
		d2, err := minifier.Bytes("text/javascript", d)
		maybePanicIfErr(err)
		if err == nil {
			lg("Minified %s from %d => %d (saved %d)\n", srcName, len(d), len(d2), len(d)-len(d2))
			d = d2
		}
	}

	d = append(book.tocData, d...)
	sha1Hex := u.Sha1HexOfBytes(d)
	name := nameToSha1Name(srcName, sha1Hex)
	dst := filepath.Join("www", "s", name)
	err = ioutil.WriteFile(dst, d, 0644)
	maybePanicIfErr(err)
	if err != nil {
		return
	}
	book.AppJSURL = "/s/" + name
	lg("Created %s\n", dst)
}
