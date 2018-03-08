package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// Article represents a part of a chapter
type Article struct {
	// stable, globally unique (across all bookd) id
	// either imported Id from Stack Overflow or auto-generated by us
	// allows stable urls and being able to cross-reference articles
	ID             string
	No             int      // TODO: can I get rid of this?
	Chapter        *Chapter // reference to containing chapter
	Title          string   // used in book_index.tmpl.html
	SearchSynonyms []string // from Search:
	FileNameBase   string   // base for both filename and url, format: ${ID}-${Title}
	BodyMarkdown   string
	// TODO: we should convert all HTML content to markdown
	BodyHTML template.HTML

	// for search we extract headings from markdown source
	cachedHeadings []HeadingInfo

	// for generating toc of a chapter, all articles that belong to the same
	// chapter as this article
	Siblings  []Article
	IsCurrent bool // only used when part of Siblings

	sourceFilePath string // path of the file from which we've read the article
}

// Book retuns book this article belongs to
func (a *Article) Book() *Book {
	return a.Chapter.Book
}

// HTML returns html content of the article
func (a *Article) HTML() template.HTML {
	if a.BodyHTML == "" {
		defLang := getDefaultLangForBook(a.Book().Title)
		html := markdownToHTML([]byte(a.BodyMarkdown), defLang, a.Book())
		a.BodyHTML = template.HTML(html)
	}
	return a.BodyHTML
}

// Headings returns headings in markdown file
func (a *Article) Headings() []HeadingInfo {
	if a.cachedHeadings != nil {
		return a.cachedHeadings
	}
	headings := parseHeadingsFromMarkdown([]byte(a.BodyMarkdown))
	a.cachedHeadings = headings
	return headings
}

// URL returns url of .html file with this article
func (a *Article) URL() string {
	chap := a.Chapter
	book := chap.Book
	// /essential/go/14047-flags
	return fmt.Sprintf("/essential/%s/%s", book.FileNameBase, a.FileNameBase)
}

// CanonnicalURL returns full url including host
func (a *Article) CanonnicalURL() string {
	return urlJoin(siteBaseURL, a.URL())
}

// GitHubText returns text we display in GitHub box
func (a *Article) GitHubText() string {
	return "Edit on GitHub"
}

// GitHubURL returns url to GitHub repo
func (a *Article) GitHubURL() string {
	uri := a.Chapter.GitHubURL() + "/" + filepath.Base(a.sourceFilePath)
	uri = strings.Replace(uri, "/tree/", "/blob/", -1)
	return uri
}

// GitHubEditURL returns url to editing this article on GitHub
// same as GitHubURL because we don't want to automatically fork
// the repo as would happen if we used /edit/ url
func (a *Article) GitHubEditURL() string {
	return a.GitHubURL()
}

// GitHubIssueURL returns link for reporting an issue about an article on githbu
// https://github.com/essentialbooks/books/issues/new?title=${title}&body=${body}&labels=docs"
func (a *Article) GitHubIssueURL() string {
	title := fmt.Sprintf("Issue for article '%s'", a.Title)
	body := fmt.Sprintf("From URL: %s\nFile: %s\n", a.CanonnicalURL(), a.GitHubEditURL())
	return gitHubBaseURL + fmt.Sprintf("/issues/new?title=%s&body=%s&labels=docs", title, body)
}

// PageTitle returns title for the page
// We want this to be unique for SEO purposes
func (a *Article) PageTitle() string {
	return a.Title + " in chapter '" + a.Chapter.Title + "'"
}

func (a *Article) destFilePath() string {
	return filepath.Join(destEssentialDir, a.Book().FileNameBase, a.FileNameBase+".html")
}
