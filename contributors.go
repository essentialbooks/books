package main

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/essentialbooks/books/pkg/common"

	"github.com/kjk/u"
)

// SoContributor describes a StackOverflow contributor
type SoContributor struct {
	ID      int
	URLPart string
	Name    string
}

func soContributorURL(userID int, userName string) string {
	return fmt.Sprintf("https://stackoverflow.com/users/%d/%s", userID, userName)
}

func loadSoContributorsMust(book *Book) {
	path := filepath.Join("books", book.Dir+"_so_contributors.txt")

	// it's ok if the file doesn't exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		lg("No contributors for book %s\n", book.Title)
		book.SoContributors = []SoContributor{}
		return
	}

	lg("loadSoContributorsMust: book.Dir: %s, path: %s\n", book.Dir, path)
	lines, err := common.ReadFileAsLines(path)
	panicIfErr(err)
	var contributors []SoContributor
	for _, line := range lines {
		id, err := strconv.Atoi(line)
		u.PanicIfErr(err)
		name := soUserIDToNameMap[id]
		u.PanicIf(name == "", "no SO contributor for id %d", id)
		if name == "user_deleted" {
			continue
		}
		nameUnescaped, err := url.PathUnescape(name)
		u.PanicIfErr(err)
		c := SoContributor{
			ID:      id,
			URLPart: name,
			Name:    nameUnescaped,
		}
		contributors = append(contributors, c)
	}
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Name < contributors[j].Name
	})
	book.SoContributors = contributors
}

func genContributorsHTML(contributors []SoContributor) string {
	if len(contributors) == 0 {
		return ""
	}
	lines := []string{
		`<div>Contributors from <a href="https://github.com/essentialbooks/books/graphs/contributors">GitHub</a>:</div>`,
		`<div></div>`,
		`<div>Contributors from Stack Overflow:</div>`,
		`<ul>`,
	}
	for _, c := range contributors {
		s := fmt.Sprintf(`<li><a href="%s">%s</a></li>`, soContributorURL(c.ID, c.Name), c.Name)
		lines = append(lines, s)
	}
	lines = append(lines, `</ul>`)
	return strings.Join(lines, "\n")
}

func genContributorsPage(book *Book) {
	loadSoContributorsMust(book)
	if book.ContributorCount() == 0 {
		return
	}
	s := genContributorsHTML(book.SoContributors)
	if s == "" {
		return
	}
	page := &Page{
		Title:    "Contributors",
		Book:     book,
		NotionID: "9999",
		BodyHTML: template.HTML(s),
	}
	book.RootPage.Pages = append(book.RootPage.Pages, page)
}
