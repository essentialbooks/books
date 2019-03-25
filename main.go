package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"html/template"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/notionapi"
	"github.com/kjk/u"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

var (
	flgAnalytics string
	flgPreview   bool
	flgAllBooks  bool
	// if true, disables downloading pages
	flgNoDownload bool
	// if true, disables notion cache, forcing re-download of notion page
	// even if cached verison on disk exits
	flgDisableNotionCache bool
	// url or id of the page to rebuild
	flgNoUpdateOutput bool

	gShowForum = true
	//gForumLink = "https://spectrum.chat/programming-books"
	gForumLink = "https://www.reddit.com/r/essentialbooks/"

	soUserIDToNameMap map[int]string
	googleAnalytics   template.HTML
	doMinify          bool
	minifier          *minify.M

	notionAuthToken string

	// when downloading pages from the server, count total number of
	// downloaded and those from cache
	nTotalDownloaded int
	nTotalFromCache  int
)

var (
	booksMain = []*Book{
		bookGo,
		bookCpp,
		bookJavaScript,
		bookCSS,
		bookHTML,
		bookHTMLCanvas,
		bookJava,
		bookKotlin,
		bookCsharp,
		bookPython,
		bookPostgresql,
		bookMysql,
		bookIOS,
		bookAndroid,
		bookBash,
		bookPowershell,
		bookBatch,
		bookGit,
		bookPHP,
		bookRuby,
		bookNode,
		bookDart,
		bookTypeScript,
		bookSwift,
	}
	booksUnpublished = []*Book{
		bookNETFramework,
		bookAlgorithm,
		bookC,
		bookObjectiveC,
		bookReact,
		bookReactNative,
		bookRubyOnRails,
		bookSql,
	}
	allBooks = append(booksMain, booksUnpublished...)
)

func parseFlags() {
	flag.StringVar(&flgAnalytics, "analytics", "", "google analytics code")
	flag.BoolVar(&flgPreview, "preview", false, "if true will start watching for file changes and re-build everything")
	flag.BoolVar(&flgAllBooks, "all-books", false, "if true will do all books")
	flag.BoolVar(&flgNoUpdateOutput, "no-update-output", false, "if true, will disable updating ouput files in cache")
	flag.BoolVar(&flgDisableNotionCache, "no-cache", false, "if true, disables cache for notion")
	flag.BoolVar(&flgNoDownload, "no-download", false, "if true, will not download pages from notion")
	flag.Parse()

	if flgAnalytics != "" {
		googleAnalyticsTmpl := `<script async src="https://www.googletagmanager.com/gtag/js?id=%s"></script>
		<script>
			window.dataLayer = window.dataLayer || [];
			function gtag(){dataLayer.push(arguments);}
			gtag('js', new Date());
			gtag('config', '%s')
		</script>
	`
		s := fmt.Sprintf(googleAnalyticsTmpl, flgAnalytics, flgAnalytics)
		googleAnalytics = template.HTML(s)
	}

	notionAuthToken = os.Getenv("NOTION_TOKEN")
	if notionAuthToken != "" {
		fmt.Printf("NOTION_TOKEN provided, can write back\n")
	} else {
		fmt.Printf("NOTION_TOKEN not provided, read only\n")
	}
}

func downloadBook(c *notionapi.Client, book *Book) {
	lg("Loading %s...", book.Title)
	pages := loadPagesFromDisk(book.NotionCacheDir())
	for _, notionPage := range pages {
		id := normalizeID(notionPage.ID)
		page := book.idToPage[id]
		if page == nil {
			page = &Page{
				NotionPage: notionPage,
			}
			book.idToPage[id] = page
		}
	}
	checkIfPagesAreOutdated(c, book.idToPage)
	loadNotionPages(c, book)
	lg("Got %d pages for %s\n", len(book.idToPage), book.Title)
	bookFromPages(book)
}

func loadSOUserMappingsMust() {
	path := filepath.Join("stack-overflow-docs-dump", "users.json.gz")
	err := common.JSONDecodeGzipped(path, &soUserIDToNameMap)
	u.PanicIfErr(err)
}

func shouldCopyImage(path string) bool {
	return !strings.Contains(path, "@2x")
}

func copyCoversMust() {
	copyFilesRecur(filepath.Join("www", "covers"), "covers", shouldCopyImage)
}

func getAlmostMaxProcs() int {
	// leave some juice for other programs
	nProcs := runtime.NumCPU() - 2
	if nProcs < 1 {
		return 1
	}
	return nProcs
}

// copy from tmpl to www, optimize if possible, add
// sha1 of the content as part of the name
func copyToWwwAsSha1MaybeMust(srcName string) {
	var dstPtr *string
	minifyType := ""
	switch srcName {
	case "main.css":
		dstPtr = &pathMainCSS
		minifyType = "text/css"
	case "app.js":
		dstPtr = &pathAppJS
		minifyType = "text/javascript"
	case "favicon.ico":
		dstPtr = &pathFaviconICO
	default:
		panicIf(true, "unknown srcName '%s'", srcName)
	}
	src := filepath.Join("tmpl", srcName)
	d, err := ioutil.ReadFile(src)
	panicIfErr(err)

	if doMinify && minifyType != "" {
		d2, err := minifier.Bytes(minifyType, d)
		maybePanicIfErr(err)
		if err == nil {
			fmt.Printf("Compressed %s from %d => %d (saved %d)\n", srcName, len(d), len(d2), len(d)-len(d2))
			d = d2
		}
	}

	sha1Hex := u.Sha1HexOfBytes(d)
	name := nameToSha1Name(srcName, sha1Hex)
	dst := filepath.Join("www", "s", name)
	err = ioutil.WriteFile(dst, d, 0644)
	panicIfErr(err)
	*dstPtr = filepath.ToSlash(dst[len("www"):])
	fmt.Printf("Copied %s => %s\n", src, dst)
}

func genBooks(books []*Book) {
	nProcs := getAlmostMaxProcs()

	timeStart := time.Now()
	clearSitemapURLS()
	copyCoversMust()

	copyToWwwAsSha1MaybeMust("main.css")
	copyToWwwAsSha1MaybeMust("app.js")
	copyToWwwAsSha1MaybeMust("favicon.ico")
	genIndex(books)
	genIndexGrid(books)
	gen404TopLevel()
	genAbout()
	genFeedback()

	for _, book := range books {
		genBook(book)
	}
	writeSitemap()
	fmt.Printf("Used %d procs, finished generating all books in %s\n", nProcs, time.Since(timeStart))
}

func initMinify() {
	minifier = minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/javascript", js.Minify)
	// less aggresive minification because html validators
	// report this as html errors
	minifier.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags:      true,
	})
	doMinify = !flgPreview
}

func isNotionCachedInDir(dir string, id string) bool {
	id = normalizeID(id)
	files, err := ioutil.ReadDir(dir)
	panicIfErr(err)
	for _, fi := range files {
		name := fi.Name()
		if strings.HasPrefix(name, id) {
			return true
		}
	}
	return false
}

func findBookFromCachedPageID(id string) *Book {
	files, err := ioutil.ReadDir("cache")
	panicIfErr(err)
	for _, book := range allBooks {
		if book.NotionStartPageID == id {
			return book
		}
	}

	for _, fi := range files {
		if !fi.IsDir() {
			continue
		}
		dir := fi.Name()
		book := findBook(dir)
		panicIf(book == nil, "didn't find book for dir '%s'", dir)
		if isNotionCachedInDir(filepath.Join("cache", dir, "notion"), id) {
			return book
		}
	}
	return nil
}

func isReplitURL(uri string) bool {
	return strings.Contains(uri, "repl.it/")
}

func initBook(book *Book) {
	var err error

	createDirMust(book.OutputCacheDir())
	createDirMust(book.NotionCacheDir())

	if false {
		loadCache("cache/go/cache.txt")
		os.Exit(0)
	}

	book.idToPage = map[string]*Page{}
	book.cache = loadCache(book.CachePath())
	must(err)
}

func findBook(id string) *Book {
	for _, book := range allBooks {
		// fuzzy match - whatever hits
		parts := []string{book.Title, book.Dir, book.NotionStartPageID}
		for _, s := range parts {
			if strings.EqualFold(s, id) {
				return book
			}
		}
	}
	return nil
}

func adHoc() {
	if false {
		glotRunTestAndExit()
	}
	if false {
		glotGetSnippedIDTestAndExit()
	}
	if false {
		// only needs to be run when we add new covers
		genTwitterImagesAndExit()
	}
}

func main() {
	openLog()
	defer closeLog()

	parseFlags()

	adHoc()

	os.RemoveAll("www")
	createDirMust(filepath.Join("www", "s"))
	createDirMust("log")

	initMinify()
	loadSOUserMappingsMust()

	client := &notionapi.Client{
		AuthToken: notionAuthToken,
	}

	books := booksMain
	if flgAllBooks {
		books = allBooks
		lg("Downloading all books\n")
	} else {
		if len(flag.Args()) > 0 {
			var newBooks []*Book
			for _, name := range flag.Args() {
				book := findBook(name)
				if book == nil {
					lg("Didn't find book named '%s'\n", name)
					continue
				}
				newBooks = append(newBooks, book)
			}
			if len(newBooks) > 0 {
				books = newBooks
				lg("Downloading %d books %#v\n", len(books), books)
			}
		}
	}

	for _, book := range books {
		initBook(book)
		downloadBook(client, book)
		loadSoContributorsMust(book)
	}

	genBooks(books)
	genNetlifyHeaders()
	genNetlifyRedirects(books)
	printAndClearErrors()

	lg("Downloaded %d pages, got %d from cache\n", nTotalDownloaded, nTotalFromCache)
	if flgPreview {
		startPreview()
	}
}
