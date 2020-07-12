package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/caching_downloader"
	"github.com/kjk/u"
)

// Book represents a book
type Book struct {
	Title     string // "Go", "jQuery" etcc
	TitleLong string // "Essential Go", "Essential jQuery" etc.

	// used by page index. defaults to: "<b>${TitleLong}</b> is a free book about ${Title} programming language."
	summary string

	NotionStartPageID string
	RootPage          *Page   // a tree of pages
	cachedPages       []*Page // pages flattened into an array

	idToPage map[string]*Page

	Dir string // directory name for the book e.g. "go"

	// generated toc javascript data
	tocData []byte
	// url of combined tocData and app.js
	AppJSURL string

	// name of a file in covers/ directory
	// e.g. "Python.png"
	CoverImageName string

	client *notionapi.Client
	// cache related
	cache *Cache
}

// CacheDir returns a cache dir for this book
func (b *Book) CacheDir() string {
	u.PanicIf(b.Dir == "", "b.Dir should not be empty")
	return filepath.Join("cache", b.Dir)
}

// NotionCacheDir returns output cache dir for this book
func (b *Book) NotionCacheDir() string {
	return filepath.Join(b.CacheDir(), "notion")
}

func (b *Book) cachePath() string {
	return filepath.Join(b.CacheDir(), "cache.txt")
}

// this is where html etc. files for a book end up
func (b *Book) destDir() string {
	return filepath.Join(destEssentialDir, b.Dir)
}

// URL returns url of the book, used in index.tmpl.html
func (b *Book) URL() string {
	return fmt.Sprintf("/essential/%s/", b.Dir)
}

func (b *Book) Summary() template.HTML {
	if b.summary == "" {
		b.summary = fmt.Sprintf("<b>%s</b> is a free book about %s programming language.", b.TitleLong, b.Title)
	}
	return template.HTML(b.summary)
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
	u.PanicIf(b.CoverImageName == "")
	return fmt.Sprintf("/covers/%s", b.CoverImageName)
}

// CoverSmallURL returns url to small cover image
func (b *Book) CoverSmallURL() string {
	u.PanicIf(b.CoverImageName == "")
	return fmt.Sprintf("/covers_small/%s", b.CoverImageName)
}

// CoverFullURL returns a URL for the cover including host
func (b *Book) CoverFullURL() string {
	return urlJoin(siteBaseURL, b.CoverURL())
}

// CoverTwitterFullURL returns a URL for the cover including host
func (b *Book) CoverTwitterFullURL() string {
	u.PanicIf(b.CoverImageName == "")
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
	d := book.tocData
	if doMinify {
		d2, err := minifier.Bytes("text/javascript", d)
		maybePanicIfErr(err)
		if err == nil {
			logf("Minified %s from %d => %d (saved %d)\n", srcName, len(d), len(d2), len(d)-len(d2))
			d = d2
		}
	}

	sha1Hex := u.Sha1HexOfBytes(d)
	name := nameToSha1Name(srcName, sha1Hex)
	dst := filepath.Join("www", "s", name)
	err := ioutil.WriteFile(dst, d, 0644)
	maybePanicIfErr(err)
	if err != nil {
		return
	}
	book.AppJSURL = "/s/" + name
	logf("Created %s\n", dst)
}

func calcPageHeadings(page *Page) {
	var headings []*HeadingInfo
	cb := func(block *notionapi.Block) {
		isHeader := false
		switch block.Type {
		case notionapi.BlockHeader, notionapi.BlockSubHeader, notionapi.BlockSubSubHeader:
			isHeader = true
		}
		if !isHeader {
			return
		}
		id := notionapi.ToNoDashID(block.ID)
		s := getInlinesPlain(block.InlineContent)
		h := &HeadingInfo{
			Text: s,
			ID:   id,
		}
		headings = append(headings, h)
	}
	blocks := []*notionapi.Block{page.NotionPage.Root()}
	notionapi.ForEachBlock(blocks, cb)
	page.Headings = headings
}

func (b *Book) afterPageDownload(page *notionapi.Page) error {
	id := toNoDashID(page.ID)
	p := &Page{
		NotionPage: page,
		NotionID:   id,
		Book:       b,
	}
	b.idToPage[id] = p
	evalCodeSnippetsForPage(p)
	downloadImages(b, p)
	calcPageHeadings(p)
	return nil
}

func initBook(book *Book) {
	var err error

	u.CreateDirMust(book.NotionCacheDir())

	if false {
		loadCache("cache/go/cache.txt")
		os.Exit(0)
	}

	book.idToPage = map[string]*Page{}
	book.cache = loadCache(book.cachePath())
	must(err)
}

func downloadBook(book *Book) {
	logf("Loading %s...\n", book.Title)
	nProcessed = 0
	nNotionPagesFromCache = 0
	nDownloadedPages = 0

	book.client = newNotionClient()
	cacheDir := book.NotionCacheDir()
	dirCache, err := caching_downloader.NewDirectoryCache(cacheDir)
	must(err)
	d := caching_downloader.New(dirCache, book.client)
	d.EventObserver = eventObserver
	d.RedownloadNewerVersions = !flgNoDownload
	d.NoReadCache = flgDisableNotionCache

	startPageID := book.NotionStartPageID
	pages, err := d.DownloadPagesRecursively(startPageID, book.afterPageDownload)
	must(err)
	nPages := len(pages)
	logf("Got %d pages for %s, downloaded: %d, from cache: %d\n", nPages, book.Title, nDownloadedPages, nNotionPagesFromCache)
	bookFromPages(book)
}
