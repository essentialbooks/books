package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	// top-level directory where .html files are generated
	destDir = "www"
	tmplDir = "tmpl"
)

var ( // directory where generated .html files for books are
	destEssentialDir = filepath.Join(destDir, "essential")
	pathAppJS        = "/s/app.js"
	pathMainCSS      = "/s/main.css"
	pathFaviconICO   = "/s/favicon.ico"
)

var (
	templateNames = []string{
		"index.tmpl.html",
		"index-grid.tmpl.html",
		"book_index.tmpl.html",
		"chapter.tmpl.html",
		"article.tmpl.html",
		"about.tmpl.html",
		"feedback.tmpl.html",
		"404.tmpl.html",
	}
	templates = make([]*template.Template, len(templateNames))

	gitHubBaseURL = "https://github.com/essentialbooks/books"
	notionBaseURL = "https://notion.so/"
	siteBaseURL   = "https://www.programming-books.io"
)

func unloadTemplates() {
	templates = make([]*template.Template, len(templateNames))
}

func tmplPath(name string) string {
	return filepath.Join(tmplDir, name)
}

var (
	funcMap = template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}
)

func loadTemplateHelperMaybeMust(name string, ref **template.Template) *template.Template {
	res := *ref
	if res != nil {
		return res
	}
	path := tmplPath(name)
	//fmt.Printf("loadTemplateHelperMust: %s\n", path)
	t, err := template.New(name).Funcs(funcMap).ParseFiles(path)
	maybePanicIfErr(err)
	if err != nil {
		return nil
	}
	*ref = t
	return t
}

func loadTemplateMaybeMust(name string) *template.Template {
	var ref **template.Template
	for i, tmplName := range templateNames {
		if tmplName == name {
			ref = &templates[i]
			break
		}
	}
	if ref == nil {
		log.Fatalf("unknown template '%s'\n", name)
	}
	return loadTemplateHelperMaybeMust(name, ref)
}

func execTemplateToFileSilentMaybeMust(name string, data interface{}, path string) error {
	var errToReturn error
	tmpl := loadTemplateMaybeMust(name)
	if tmpl == nil {
		return nil
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	maybePanicIfErr(err)

	d := buf.Bytes()
	if doMinify {
		d2, err := minifier.Bytes("text/html", d)
		//maybePanicIfErr(err)
		if err == nil {
			totalHTMLBytes += len(d)
			totalHTMLBytesMinified += len(d2)
			d = d2
		} else {
			errToReturn = err
		}
	}
	err = ioutil.WriteFile(path, d, 0644)
	maybePanicIfErr(err)
	return errToReturn
}

func execTemplateToFileMaybeMust(name string, data interface{}, path string) error {
	return execTemplateToFileSilentMaybeMust(name, data, path)
}

// PageCommon is a common information for most pages
type PageCommon struct {
	Analytics      template.HTML
	PathAppJS      string
	PathMainCSS    string
	PathFaviconICO string
}

func getPageCommon() PageCommon {
	return PageCommon{
		Analytics:      googleAnalytics,
		PathAppJS:      pathAppJS,
		PathMainCSS:    pathMainCSS,
		PathFaviconICO: pathFaviconICO,
	}
}

func gen404TopLevel() {
	d := struct {
		PageCommon
		Book *Book
	}{
		PageCommon: getPageCommon(),
	}
	path := filepath.Join(destDir, "404.html")
	execTemplateToFileMaybeMust("404.tmpl.html", d, path)
}

func genIndex(books []*Book) {
	d := struct {
		PageCommon
		Books     []*Book
		NotionURL string
	}{
		PageCommon: getPageCommon(),
		Books:      books,
		NotionURL:  gitHubBaseURL,
	}
	path := filepath.Join(destDir, "index.html")
	execTemplateToFileMaybeMust("index.tmpl.html", d, path)
}

func genIndexGrid(books []*Book) {
	d := struct {
		PageCommon
		Books []*Book
	}{
		PageCommon: getPageCommon(),
		Books:      books,
	}
	path := filepath.Join(destDir, "index-grid.html")
	execTemplateToFileMaybeMust("index-grid.tmpl.html", d, path)
}

func genFeedback() {
	fmt.Printf("writing feedback.html\n")
	path := filepath.Join(destDir, "feedback.html")
	d := struct {
		PageCommon
		ForumLink string
	}{
		PageCommon: getPageCommon(),
		ForumLink:  gForumLink,
	}
	execTemplateToFileMaybeMust("feedback.tmpl.html", d, path)
}

func genAbout() {
	d := getPageCommon()
	fmt.Printf("writing about.html\n")
	path := filepath.Join(destDir, "about.html")
	execTemplateToFileMaybeMust("about.tmpl.html", d, path)
}

// TODO: consolidate chapter/article html
func genArticle(book *Book, page *Page, currChapNo int, currArticleNo int) {
	addSitemapURL(page.CanonnicalURL())

	d := struct {
		PageCommon
		*Page
		CurrentChapterNo int
		CurrentArticleNo int
		ShowForum        bool
		ForumLink        string
	}{
		PageCommon:       getPageCommon(),
		Page:             page,
		CurrentChapterNo: currChapNo,
		CurrentArticleNo: currArticleNo,
		ShowForum:        gShowForum,
		ForumLink:        gForumLink,
	}

	path := page.destFilePath()
	err := execTemplateToFileSilentMaybeMust("article.tmpl.html", d, path)
	if err != nil {
		fmt.Printf("Failed to minify page %s in book %s\n", page.NotionID, book.Title)
	}
}

func genChapter(book *Book, page *Page, currNo int) {
	addSitemapURL(page.CanonnicalURL())
	for i, article := range page.Pages {
		genArticle(book, article, currNo, i)
	}

	path := page.destFilePath()
	d := struct {
		PageCommon
		*Page
		CurrentChapterNo int
		ShowForum        bool
		ForumLink        string
	}{
		PageCommon:       getPageCommon(),
		Page:             page,
		CurrentChapterNo: currNo,
		ShowForum:        gShowForum,
		ForumLink:        gForumLink,
	}
	execTemplateToFileSilentMaybeMust("chapter.tmpl.html", d, path)

	for _, imagePath := range page.images {
		imageName := filepath.Base(imagePath)
		dst := page.destImagePath(imageName)
		copyFileMaybeMust(dst, imagePath)
	}
}

func buildIDToPage(book *Book) {
	pages := book.GetAllPages()
	for _, page := range pages {
		id := normalizeID(page.NotionPage.ID)
		book.idToPage[id] = page
		page.Book = book
	}
}

func bookPagesToHTML(book *Book) {
	nProcessed := 0
	pages := book.GetAllPages()
	for _, page := range pages {
		html := notionToHTML(page, book)
		page.BodyHTML = template.HTML(string(html))
		nProcessed++
	}
	fmt.Printf("bookPagesToHTML: processed %d pages for book %s\n", nProcessed, book.TitleLong)
}

func bookPageToHTML(book *Book, id string) {
	pages := book.GetAllPages()
	for _, page := range pages {
		if page.NotionID == id {
			fmt.Printf("bookPageToHTML: processed page %s for book %s\n", id, book.TitleLong)
			html := notionToHTML(page, book)
			page.BodyHTML = template.HTML(string(html))
			return
		}
	}
	fmt.Printf("bookPageToHTML: didn't find page '%s' for book %s\n", id, book.TitleLong)
}

func genBook(book *Book) {
	lg("Started genering book %s\n", book.Title)
	timeStart := time.Now()

	buildIDToPage(book)
	genContributorsPage(book)
	bookPagesToHTML(book)

	genBookTOCSearchMust(book)

	// generate index.html for the book
	err := os.MkdirAll(book.destDir(), 0755)
	maybePanicIfErr(err)
	if err != nil {
		return
	}

	data := struct {
		PageCommon
		Book *Book
	}{
		PageCommon: getPageCommon(),
		Book:       book,
	}

	path := filepath.Join(book.destDir(), "index.html")
	execTemplateToFileSilentMaybeMust("book_index.tmpl.html", data, path)

	// TODO: per-book 404 should link to top of book, not top of website
	path = filepath.Join(book.destDir(), "404.html")
	execTemplateToFileSilentMaybeMust("404.tmpl.html", data, path)
	addSitemapURL(book.CanonnicalURL())

	for i, chapter := range book.Chapters() {
		genChapter(book, chapter, i)
	}

	fmt.Printf("Generated book '%s' in %s\n", book.Title, time.Since(timeStart))
}
