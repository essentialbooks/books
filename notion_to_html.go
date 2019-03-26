package main

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"

	"github.com/alecthomas/template"
	"github.com/kjk/notionapi"
)

// HTMLGenerator is for notion -> HTML generation
type HTMLGenerator struct {
	f            *bytes.Buffer
	page         *Page
	level        int
	nToggle      int
	err          error
	book         *Book
	currHeaderID int
}

// only hex chars seem to be valid
func isValidNotionIDChar(c byte) bool {
	switch {
	case c >= '0' && c <= '9':
		return true
	case c >= 'a' && c <= 'f':
		return true
	case c >= 'A' && c <= 'F':
		// currently not used but just in case they change their minds
		return true
	}
	return false
}

func isValidNotionID(id string) bool {
	// len("ea07db1b9bff415ab180b0525f3898f6")
	if len(id) != 32 {
		return false
	}
	for i := range id {
		if !isValidNotionIDChar(id[i]) {
			return false
		}
	}
	return true
}

// https://www.notion.so/Advanced-web-spidering-with-Puppeteer-ea07db1b9bff415ab180b0525f3898f6
// https://www.notion.so/c674bebe8adf44d18c3a36cc18c131e2
// returns "" if didn't detect valid notion id in the url
func extractNotionIDFromURL(uri string) string {
	if isValidNotionID(normalizeID(uri)) {
		return uri
	}
	trimmed := strings.TrimPrefix(uri, "https://www.notion.so/")
	if uri == trimmed {
		return ""
	}
	// could be c674bebe8adf44d18c3a36cc18c131e2 from https://www.notion.so/c674bebe8adf44d18c3a36cc18c131e2
	id := trimmed
	parts := strings.Split(trimmed, "-")
	n := len(parts)
	if n >= 2 {
		// could be ea07db1b9bff415ab180b0525f3898f6 from Advanced-web-spidering-with-Puppeteer-ea07db1b9bff415ab180b0525f3898f6
		id = parts[n-1]
	}
	id = normalizeID(id)
	if !isValidNotionID(id) {
		return ""
	}
	return id
}

// change https://www.notion.so/Advanced-web-spidering-with-Puppeteer-ea07db1b9bff415ab180b0525f3898f6
// =>
// /article/${id}
func (g *HTMLGenerator) maybeReplaceNotionLink(uri string) string {
	id := extractNotionIDFromURL(uri)
	if id == "" {
		return uri
	}
	page := g.book.idToPage[id]
	if page == nil {
		lg("Didn't find page with id '%s' extracted from url %s\n", id, uri)
		return uri
	}
	return page.URL()
}

func (g *HTMLGenerator) getURLAndTitleForBlock(block *notionapi.Block) (string, string) {
	id := normalizeID(block.ID)
	page := g.book.idToPage[id]
	if page == nil {
		title := cleanTitle(block.Title)
		lg("No article for id %s %s\n", id, title)
		url := "/article/" + id + "/" + urlify(title)
		return url, title
	}

	return page.URL(), page.Title
}

func findPageByID(book *Book, id string) *Page {
	pages := book.GetAllPages()
	for _, page := range pages {
		if strings.EqualFold(page.getID(), id) {
			return page
		}
	}
	return nil
}

func (g *HTMLGenerator) reportIfInvalidLink(uri string) {
	link := g.maybeReplaceNotionLink(uri)
	if link != uri {
		return
	}
	if strings.HasPrefix(uri, "http") {
		return
	}
	pageID := normalizeID(g.page.getID())
	lg("Found invalid link '%s' in page https://notion.so/%s", uri, pageID)
	destPage := findPageByID(g.book, uri)
	if destPage != nil {
		lg(" most likely pointing to https://notion.so/%s\n", normalizeID(destPage.NotionPage.ID))
	} else {
		lg("\n")
	}
}

func (g *HTMLGenerator) genInlineBlock(b *notionapi.InlineBlock) {
	var start, close string
	if b.AttrFlags&notionapi.AttrBold != 0 {
		start += "<b>"
		close += "</b>"
	}
	if b.AttrFlags&notionapi.AttrItalic != 0 {
		start += "<i>"
		close += "</i>"
	}
	if b.AttrFlags&notionapi.AttrStrikeThrought != 0 {
		start += "<strike>"
		close += "</strike>"
	}
	if b.AttrFlags&notionapi.AttrCode != 0 {
		start += "<code>"
		close += "</code>"
	}
	skipText := false
	if b.Link != "" {
		g.reportIfInvalidLink(b.Link)
		link := g.maybeReplaceNotionLink(b.Link)
		start += fmt.Sprintf(`<a href="%s">%s</a>`, link, b.Text)
		skipText = true
	}
	if b.UserID != "" {
		start += fmt.Sprintf(`<span class="user">@%s</span>`, b.UserID)
		skipText = true
	}
	if b.Date != nil {
		// TODO: serialize date properly
		start += fmt.Sprintf(`<span class="date">@TODO: date</span>`)
		skipText = true
	}
	if !skipText {
		start += b.Text
	}
	g.writeString(start + close)
}

func (g *HTMLGenerator) getInline(blocks []*notionapi.InlineBlock) []byte {
	b := g.newBuffer()
	g.genInlineBlocks(blocks)
	return g.restoreBuffer(b)
}

func (g *HTMLGenerator) genInlineBlocks(blocks []*notionapi.InlineBlock) {
	for _, block := range blocks {
		g.genInlineBlock(block)
	}
}

func genInlineBlocksText(blocks []*notionapi.InlineBlock) string {
	var a []string
	for _, b := range blocks {
		a = append(a, b.Text)
	}
	return strings.Join(a, "")
}

// In notion I want to have @TODO lines that are not rendered in html output
func isTodoBlock(block *notionapi.Block) bool {
	panicIf(block.Type != notionapi.BlockText, "only supported on '%s' block, called on '%s' block", notionapi.BlockText, block.Type)
	blocks := block.InlineContent
	if len(blocks) == 0 {
		return false
	}
	b := blocks[0]
	if strings.HasPrefix(b.Text, "@TODO") {
		return true
	}
	return false
}

func (g *HTMLGenerator) genBlockSurrouded(block *notionapi.Block, start, close string) {
	g.writeString(start + "\n")
	g.genInlineBlocks(block.InlineContent)
	g.level++
	g.genContent(block)
	g.level--
	g.writeString(close + "\n")
}

/*
v is expected to be
[
	[
		"foo"
	]
]
and we want to return "foo"
If not present or unexpected shape, return ""
is still visible
*/
func propsValueToText(v interface{}) string {
	if v == nil {
		return ""
	}

	// [ [ "foo" ]]
	a, ok := v.([]interface{})
	if !ok {
		return fmt.Sprintf("type1: %T", v)
	}
	// [ "foo" ]
	if len(a) == 0 {
		return ""
	}
	v = a[0]
	a, ok = v.([]interface{})
	if !ok {
		return fmt.Sprintf("type2: %T", v)
	}
	// "foo"
	if len(a) == 0 {
		return ""
	}
	v = a[0]
	str, ok := v.(string)
	if !ok {
		return fmt.Sprintf("type3: %T", v)
	}
	return str
}

func (g *HTMLGenerator) genEmbed(block *notionapi.Block) {
	uri := block.FormatEmbed.DisplaySource
	if strings.Contains(uri, "onlinetool.io/") {
		g.genGitEmbed(block)
		return
	}
	if strings.Contains(uri, "repl.it/") {
		g.genReplitEmbed(block)
		return
	}
	panicIf(true, "unsupported embed %s", uri)
}

func (g *HTMLGenerator) genReplitEmbed(block *notionapi.Block) {
	uri := block.FormatEmbed.DisplaySource
	uri = strings.Replace(uri, "?lite=true", "", -1)
	lg("Page: https://notion.so/%s\n", g.page.NotionID)
	lg("  Replit: %s\n", uri)
	panic("we no longer use replit")
}

func (g *HTMLGenerator) genSourceFile(f *SourceFile) {
	{
		var tmp bytes.Buffer
		code := f.CodeToShow()
		lang := f.Lang
		htmlHighlight(&tmp, string(code), lang, "")
		d := tmp.Bytes()
		info := CodeBlockInfo{
			Lang:      f.Lang,
			GitHubURI: f.GitHubURL,
		}
		info.PlaygroundURI = f.PlaygroundURI
		s := fixupHTMLCodeBlock(string(d), &info)
		g.f.WriteString(s)
	}

	output := f.Output()
	if len(output) != 0 {
		var tmp bytes.Buffer
		htmlHighlight(&tmp, output, "text", "")
		d := tmp.Bytes()
		info := CodeBlockInfo{
			Lang: "output",
		}
		s := fixupHTMLCodeBlock(string(d), &info)
		g.f.WriteString(s)
	}
}

func (g *HTMLGenerator) genGitEmbed(block *notionapi.Block) {
	uri := block.FormatEmbed.DisplaySource
	f := findSourceFileForEmbedURL(g.page, uri)
	// currently we only handle source code file embeds but might handle
	// others (graphs etc.)
	if f == nil {
		lg("genEmbed: didn't find source file for url %s\n", uri)
		return
	}

	g.genSourceFile(f)
}

func (g *HTMLGenerator) genCollectionView(block *notionapi.Block) {
	viewInfo := block.CollectionViews[0]
	view := viewInfo.CollectionView
	if view.Format == nil {
		lg("genCollectionView: missing view.Format block id: %s\n", block.ID)
		return
	}

	s := `<table class="notion-table">`

	columns := view.Format.TableProperties
	s += `<thead><tr>`
	for _, col := range columns {
		colName := col.Property
		colInfo := viewInfo.Collection.CollectionSchema[colName]
		name := ""
		if colInfo != nil {
			name = colInfo.Name
		} else {
			lg("Missing colInfo in block ID '%s', page: %s\n", block.ID, g.page.NotionID)
		}
		s += `<th>` + html.EscapeString(name) + `</th>`
	}
	s += `</tr></thead>`

	s += `<tbody>`
	for _, row := range viewInfo.CollectionRows {
		s += `<tr>`
		props := row.Properties
		for _, col := range columns {
			colName := col.Property
			v := props[colName]
			colVal := propsValueToText(v)
			if colVal == "" {
				// use &nbsp; so that empty row still shows up
				// could also set a min-height to 1em or sth. like that
				s += `<td>&nbsp;</td>`
			} else {
				//colInfo := viewInfo.Collection.CollectionSchema[colName]
				// TODO: format colVal according to colInfo
				s += `<td>` + html.EscapeString(colVal) + `</td>`
			}
		}
		s += `</tr>`
	}
	s += `</tbody>`
	s += `</table>`
	g.writeString(s)
}

// Children of BlockColumnList are BlockColumn blocks
func (g *HTMLGenerator) genColumnList(block *notionapi.Block) {
	panicIf(block.Type != notionapi.BlockColumnList, "unexpected block type '%s'", block.Type)
	nColumns := len(block.Content)
	panicIf(nColumns == 0, "has no columns")
	// TODO: for now equal width columns
	s := `<div class="column-list">`
	g.writeString(s)

	for _, col := range block.Content {
		// TODO: get column ration from col.FormatColumn.ColumnRation, which is float 0...1
		panicIf(col.Type != notionapi.BlockColumn, "unexpected block type '%s'", col.Type)
		g.writeString(`<div>`)
		g.genBlocks(col.Content)
		g.writeString(`</div>`)
	}

	s = `</div>`
	g.writeString(s)
}

func (g *HTMLGenerator) newBuffer() *bytes.Buffer {
	curr := g.f
	g.f = &bytes.Buffer{}
	return curr
}

func (g *HTMLGenerator) restoreBuffer(b *bytes.Buffer) []byte {
	d := g.f.Bytes()
	g.f = b
	return d
}

func (g *HTMLGenerator) genToggle(block *notionapi.Block) {
	panicIf(block.Type != notionapi.BlockToggle, "unexpected block type '%s'", block.Type)
	g.nToggle++
	id := strconv.Itoa(g.nToggle)

	inline := g.getInline(block.InlineContent)

	b := g.newBuffer()
	g.genBlocks(block.Content)
	inner := g.restoreBuffer(b)

	s := fmt.Sprintf(`<div style="width: 100%%; margin-top: 2px; margin-bottom: 1px;">
    <div style="display: flex; align-items: flex-start; width: 100%%; padding-left: 2px; color: rgb(66, 66, 65);">

        <div style="margin-right: 4px; width: 24px; flex-grow: 0; flex-shrink: 0; display: flex; align-items: center; justify-content: center; min-height: calc((1.5em + 3px) + 3px); padding-right: 2px;">
            <div id="toggle-toggle-%s" onclick="javascript:onToggleClick(this)" class="toggler" style="align-items: center; user-select: none; display: flex; width: 1.25rem; height: 1.25rem; justify-content: center; flex-shrink: 0;">

                <svg id="toggle-closer-%s" width="100%%" height="100%%" viewBox="0 0 100 100" style="fill: currentcolor; display: none; width: 0.6875em; height: 0.6875em; transition: transform 300ms ease-in-out; transform: rotateZ(180deg);">
                    <polygon points="5.9,88.2 50,11.8 94.1,88.2 "></polygon>
                </svg>

                <svg id="toggle-opener-%s" width="100%%" height="100%%" viewBox="0 0 100 100" style="fill: currentcolor; display: block; width: 0.6875em; height: 0.6875em; transition: transform 300ms ease-in-out; transform: rotateZ(90deg);">
                    <polygon points="5.9,88.2 50,11.8 94.1,88.2 "></polygon>
                </svg>
            </div>
        </div>

        <div style="flex: 1 1 0px; min-width: 1px;">
            <div style="display: flex;">
                <div style="padding-top: 3px; padding-bottom: 3px">%s</div>
            </div>

            <div style="margin-left: -2px; display: none" id="toggle-content-%s">
                <div style="display: flex; flex-direction: column;">
                    <div style="width: 100%%; margin-top: 2px; margin-bottom: 0px;">
                        <div style="color: rgb(66, 66, 65);">
							<div style="">
								%s
                                <!-- <div style="padding: 3px 2px;">text inside list</div> -->
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
`, id, id, id, string(inline), id, string(inner))
	g.writeString(s)
}

func (g *HTMLGenerator) writeString(s string) {
	io.WriteString(g.f, s)
}

func (g *HTMLGenerator) genBlock(block *notionapi.Block) {
	levelCls := ""
	if g.level > 0 {
		levelCls = fmt.Sprintf(" lvl%d", g.level)
	}

	switch block.Type {
	case notionapi.BlockText:
		if isTodoBlock(block) {
			break
		}
		start := `<p>`
		close := `</p>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockHeader:
		g.currHeaderID++
		h := HeadingInfo{
			Text: genInlineBlocksText(block.InlineContent),
			// TODO: ID should be url-ified text
			ID: strconv.Itoa(g.currHeaderID),
		}
		g.page.Headings = append(g.page.Headings, h)
		start := fmt.Sprintf(`<h1 class="hdr%s" id="%s">`, levelCls, h.ID)
		close := `</h1>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockSubHeader:
		g.currHeaderID++
		h := HeadingInfo{
			Text: genInlineBlocksText(block.InlineContent),
			// TODO: ID should be url-ified text
			ID: strconv.Itoa(g.currHeaderID),
		}
		g.page.Headings = append(g.page.Headings, h)
		start := fmt.Sprintf(`<h2 class="hdr%s" id="%s">`, levelCls, h.ID)
		close := `</h2>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockSubSubHeader:
		g.currHeaderID++
		h := HeadingInfo{
			Text: genInlineBlocksText(block.InlineContent),
			// TODO: ID should be url-ified text
			ID: strconv.Itoa(g.currHeaderID),
		}
		g.page.Headings = append(g.page.Headings, h)
		start := fmt.Sprintf(`<h3 class="hdr%s" id="%s">`, levelCls, h.ID)
		close := `</h3>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockTodo:
		clsChecked := ""
		if block.IsChecked {
			clsChecked = " todo-checked"
		}
		start := fmt.Sprintf(`<div class="todo%s%s">`, levelCls, clsChecked)
		close := `</div>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockToggle:
		g.genToggle(block)
	case notionapi.BlockQuote:
		start := fmt.Sprintf(`<blockquote class="%s">`, levelCls)
		close := `</blockquote>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockDivider:
		fmt.Fprintf(g.f, `<hr class="%s"/>`+"\n", levelCls)
	case notionapi.BlockPage:
		cls := "page"
		if block.IsLinkToPage() {
			cls = "page-link"
		}
		url, title := g.getURLAndTitleForBlock(block)
		title = template.HTMLEscapeString(title)
		html := fmt.Sprintf(`<div class="%s%s"><a href="%s">%s</a></div>`, cls, levelCls, url, title)
		fmt.Fprintf(g.f, "%s\n", html)
	case notionapi.BlockCode:
		//lang := getLangFromFileExt(filepath.Ext(path))
		//gitHubURL := getGitHubPathForFile(path)
		lang := block.CodeLanguage
		sf := &SourceFile{
			NotionOriginURL: fmt.Sprintf("https://notion.so/%s", normalizeID(g.page.NotionID)),
			//Path:      path,
			//FileName:  name,
			//GitHubURL: gitHubURL,
		}
		sf.Lang = lang
		sf.SnippetName = g.page.PageTitle()
		if sf.SnippetName == "" {
			sf.SnippetName = "untitled"
		}

		data := []byte(block.Code)
		err := setSourceFileData(sf, data)
		if err != nil {
			lg("genBlock: setSourceFileData() failed with '%s'\n", err)
			lg("page: %s\n", sf.NotionOriginURL)
			//panicIfErr(err)
		}

		if sf.Directive.Glot || sf.Directive.GoPlayground {
			// for those we respect no output/no playground
		} else {
			// for embedded code blocks by default we don't set playground
			// or output unless explicitly asked for
			sf.Directive.NoPlayground = true
			sf.Directive.NoOutput = true
		}
		setDefaultFileNameFromLanguage(sf)
		err = getOutputCached(g.book.cache, sf)
		if err != nil {
			lg("getOutputCached() failed.\nsf.CodeToRun():\n%s\n", sf.CodeToRun)
			panicIfErr(err)
		}
		g.genSourceFile(sf)

		if false {
			/*
				code := template.HTMLEscapeString(block.Code)
				fmt.Fprintf(g.f, `<div class="%s">Lang for code: %s</div>
				<pre class="%s">
				%s
				</pre>`, levelCls, block.CodeLanguage, levelCls, code)
			*/
			var tmp bytes.Buffer
			htmlHighlight(&tmp, string(block.Code), block.CodeLanguage, "")
			d := tmp.Bytes()
			var info CodeBlockInfo
			// TODO: set Lang, GitHubURI and PlaygroundURI
			s := fixupHTMLCodeBlock(string(d), &info)
			g.f.WriteString(s)
		}
	case notionapi.BlockBookmark:
		fmt.Fprintf(g.f, `<div class="bookmark %s">Bookmark to %s</div>`+"\n", levelCls, block.Link)
	case notionapi.BlockGist:
		s := fmt.Sprintf(`<script src="%s.js"></script>`, block.Source)
		g.writeString(s)
	case notionapi.BlockImage:
		link := block.ImageURL
		cls := levelCls
		if len(cls) > 0 {
			cls += " "
		}
		cls += "img"
		fmt.Fprintf(g.f, `<img class="%s" src="%s" />`+"\n", cls, link)
	case notionapi.BlockColumnList:
		g.genColumnList(block)
	case notionapi.BlockCollectionView:
		g.genCollectionView(block)
	case notionapi.BlockEmbed:
		g.genEmbed(block)
	default:
		uri := "https://notion.so/" + normalizeID(g.page.NotionID)
		fmt.Printf("Unsupported block type '%s', id: %s in page %s\n", block.Type, block.ID, uri)
		panic(fmt.Sprintf("Unsupported block type '%s'", block.Type))
	}
}

func setDefaultFileNameFromLanguage(sf *SourceFile) error {
	if sf.Directive.FileName != "" {
		return nil
	}

	// we don't care unless it goes to glot.io
	if !sf.Directive.Glot {
		return nil
	}

	ext := ""
	lang := strings.ToLower(sf.Lang)
	switch lang {
	case "go":
		ext = ".go"
	case "javascript":
		ext = ".js"
	case "cpp", "cplusplus", "c++":
		ext = ".cpp"
	default:
		fmt.Printf("detectFileNameFromLanguage: lang '%s' is not supported\n", sf.Lang)
		fmt.Printf("Notion page: %s\n", sf.NotionOriginURL)
		panic("")
	}
	sf.Directive.FileName = "main" + ext
	if sf.FileName == "" {
		sf.FileName = sf.Directive.FileName
		sf.Path = sf.FileName
	}
	return nil
}

func (g *HTMLGenerator) genBlocks(blocks []*notionapi.Block) {
	for len(blocks) > 0 {
		block := blocks[0]
		if block == nil {
			lg("Missing block\n")
			blocks = blocks[1:]
			continue
		}

		if block.Type == notionapi.BlockNumberedList {
			fmt.Fprintf(g.f, `<ol>`)
			for len(blocks) > 0 {
				block := blocks[0]
				if block.Type != notionapi.BlockNumberedList {
					break
				}
				g.genBlockSurrouded(block, `<li>`, `</li>`)
				blocks = blocks[1:]
			}
			fmt.Fprintf(g.f, `</ol>`)
		} else if block.Type == notionapi.BlockBulletedList {
			fmt.Fprintf(g.f, `<ul>`)
			for len(blocks) > 0 {
				block := blocks[0]
				if block.Type != notionapi.BlockBulletedList {
					break
				}
				g.genBlockSurrouded(block, `<li>`, `</li>`)
				blocks = blocks[1:]
			}
			fmt.Fprintf(g.f, `</ul>`)
		} else {
			g.genBlock(block)
			blocks = blocks[1:]
		}
	}
}

func (g *HTMLGenerator) genContent(parent *notionapi.Block) {
	g.genBlocks(parent.Content)
}

// Gen returns generated HTML
func (g *HTMLGenerator) Gen() []byte {
	rootPage := g.page.NotionPage.Root
	f := rootPage.FormatPage
	g.writeString(`<p></p>`)
	if f != nil && f.PageFont == "mono" {
		g.writeString(`<div style="font-family: monospace">`)
	}
	g.genContent(rootPage)
	if f != nil && f.PageFont == "mono" {
		g.writeString(`</div>`)
	}
	return g.f.Bytes()
}

func notionToHTML(page *Page, book *Book) []byte {
	// This is artificially generated page (e.g. contributors page)
	if page.NotionPage == nil {
		return []byte(page.BodyHTML)
	}

	verbose("Generating HTML for %s\n", page.NotionURL())
	gen := HTMLGenerator{
		f:    &bytes.Buffer{},
		book: book,
		page: page,
	}
	return gen.Gen()
}
