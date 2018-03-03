package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/kjk/u"
)

/*
Generates a javascript file that looks like:

gBookToc = [
	[${is_expanded}, ${chapter or aticle id}, ${parentIdx}, ${childIdx}, ${title}, ${synonym 1}, ${synonym 2}, ...],
];

It's saved in wwww/essential/${bookname}/toc_search.js

Maybe: optimize this as a one long string to search instead of multiple strings
in an array. Use 0 to separate chapter/article. Use 1 to separate title
from synonims.

Also, have original and lower-cased version of the string. We search lower-cased
but show the original. That avoids lowercasing during search.
*/

const (
	itemIdxIsExpanded   = 0
	itemIdxURL          = 1
	itemIdxParent       = 2
	itemIdxFirstChild   = 3
	itemIdxTitle        = 4
	itemIdxFirstSynonym = 5
)

func genBookTOCSearchMust(book *Book) {
	var toc [][]interface{}
	for _, chapter := range book.Chapters {
		title := strings.TrimSpace(chapter.Title)
		uri := chapter.FileNameBase
		tocItem := []interface{}{false, uri, -1, -1, title}
		toc = append(toc, tocItem)
		chapIdx := len(toc) - 1
		u.PanicIf(chapIdx < 0)

		headings := chapter.Headings()
		for _, heading := range headings {
			title := heading.Text
			id := heading.ID
			if len(id) > 0 {
				id = uri + "#" + id
			}
			tocItem = []interface{}{false, id, chapIdx, -1, title}
			toc = append(toc, tocItem)
		}

		for _, article := range chapter.Articles {
			title := strings.TrimSpace(article.Title)
			uri := article.FileNameBase
			tocItem = []interface{}{false, uri, chapIdx, -1, title}
			for _, syn := range article.SearchSynonyms {
				tocItem = append(tocItem, syn)
			}
			toc = append(toc, tocItem)

			headings := article.Headings()
			articleIdx := len(toc) - 1
			for _, heading := range headings {
				title := heading.Text
				id := heading.ID
				if len(id) > 0 {
					id = uri + "#" + id
				}
				tocItem = []interface{}{false, id, articleIdx, -1, title}
				toc = append(toc, tocItem)
			}
		}
	}

	// set first child idx from parent idx
	for i, tocItem := range toc {
		parentIdx := tocItem[itemIdxParent].(int)
		if parentIdx != -1 {
			parentTocItem := toc[parentIdx]
			idx := parentTocItem[itemIdxFirstChild].(int)
			if idx == -1 {
				parentTocItem[itemIdxFirstChild] = i
			}
		}
	}

	d, err := json.MarshalIndent(toc, "", "  ")
	u.PanicIfErr(err)
	s := "gBookToc = " + string(d) + ";"
	d = []byte(s)
	path := filepath.Join(destEssentialDir, book.FileNameBase, "toc_search.js")
	u.CreateDirForFile(path)
	err = ioutil.WriteFile(path, d, 0644)
	u.PanicIfErr(err)
	fmt.Printf("%s\n", path)
}
