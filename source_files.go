package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/notionapi"
)

/*
FileDirective describes reulst of parsing a line like:
// no output, no playground
*/
type FileDirective struct {
	FileName     string // :file foo.txt
	NoOutput     bool   // "no output"
	AllowError   bool   // "allow error"
	LineLimit    int    // limit ${n}
	NoPlayground bool   // no playground
	RunCmd       string // :run ${cmd}

	Glot         bool // :glot, use glot.io to execute the code snippet
	GoPlayground bool // :goplay, use go playground to execute the snippet
	DoOutput     bool // :output
}

// SourceFile represents source file. It comes from a code block in
// Notion, replit, file in repository etc.
type SourceFile struct {
	// for debugging, url of page on notion from which this
	// snippet comes from
	NotionOriginURL string

	EmbedURL string

	// full path of the file
	Path string
	// name of the file
	FileName string

	SnippetName string

	// URL on GitHub for this file
	GitHubURL string
	// language of the file, detected from name
	Lang string

	// for Go files, this is playground id
	GoPlaygroundID string
	// for some files, this is glot.io snippet id
	GlotPlaygroundID string

	PlaygroundURI string

	// optional, extracted from first line of the file
	// allows providing meta-data instruction for this file
	Directive *FileDirective

	// raw content of the code snippet with line endings normalized to '\n'
	// it can either come from code block in Notion or a file on disk
	// or replit etc.
	CodeSnippet []byte

	LinesRaw []string // Data split into lines

	// LinesRaw after extracting directive, run cmd at the top
	// and removing :show annotation lines
	// This is the content to execute
	LinesToRun []string

	// the part that we want to show i.e. the parts inside
	// :show start, :show end blocks
	LinesCode []string

	// output of running a file
	Output string
}

// DataToRun returns content of the file after filtering, that's the
// version we want to execute to get the output
func (f *SourceFile) DataToRun() []byte {
	s := strings.Join(f.LinesToRun, "\n")
	return []byte(s)
}

// DataCode returns part of the file tbat we want to show
func (f *SourceFile) DataCode() []byte {
	s := strings.Join(f.LinesCode, "\n")
	return []byte(s)
}

// strip "//" or "#" comment mark from line and return string
// after removing the mark
func stripComment(line string) (string, bool) {
	line = strings.TrimSpace(line)
	s := strings.TrimPrefix(line, "//")
	if s != line {
		return s, true
	}
	s = strings.TrimPrefix(line, "#")
	if s != line {
		return s, true
	}
	return "", false
}

/* Parses a line like:
// no output, no playground, line ${n}, allow error
*/
func parseFileDirective(res *FileDirective, line string) (bool, error) {
	s, ok := stripComment(line)
	if !ok {
		// doesn't start with a comment, so is not a file directive
		return false, nil
	}
	parts := strings.Split(s, ",")
	for _, s := range parts {
		s = strings.TrimSpace(s)
		// directives can also start with ":", to make them more distinct
		startsWithColon := strings.HasPrefix(s, ":")
		s = strings.TrimPrefix(s, ":")
		if s == "glot" {
			res.Glot = true
		} else if s == "output" {
			res.DoOutput = true
		} else if s == "goplay" {
			res.GoPlayground = true
		} else if s == "no output" || s == "nooutput" {
			res.NoOutput = true
		} else if s == "no playground" || s == "noplayground" {
			res.NoPlayground = true
		} else if s == "allow error" || s == "allow_error" {
			res.AllowError = true
		} else if strings.HasPrefix(s, "name ") {
			// expect: name foo.txt
			rest := strings.TrimSpace(strings.TrimPrefix(s, "name "))
			if len(rest) == 0 {
				return false, fmt.Errorf("parseFileDirective: invalid line '%s'", line)
			}
			res.FileName = rest
		} else if strings.HasPrefix(s, "file ") {
			// expect: file foo.txt
			rest := strings.TrimSpace(strings.TrimPrefix(s, "file "))
			if len(rest) == 0 {
				return false, fmt.Errorf("parseFileDirective: invalid line '%s'", line)
			}
			res.FileName = rest
		} else if strings.HasPrefix(s, "line ") {
			rest := strings.TrimSpace(strings.TrimPrefix(s, "line "))
			n, err := strconv.Atoi(rest)
			if err != nil {
				return false, fmt.Errorf("parseFileDirective: invalid line '%s'", line)
			}
			res.LineLimit = n
		} else if strings.HasPrefix(s, "run ") {
			rest := strings.TrimSpace(strings.TrimPrefix(s, "run "))
			res.RunCmd = rest
			// fmt.Printf("  run:: '%s'\n", res.RunCmd)
		} else {
			// if started with ":" we assume it was meant to be a directive
			// but there was a typo
			if startsWithColon {
				return false, fmt.Errorf("parseFileDirective: invalid line '%s'", line)
			}
			// otherwise we assume this is just a comment
			return false, nil
		}
	}
	return true, nil
}

func extractFileDirectives(lines []string) (*FileDirective, []string, error) {
	directive := &FileDirective{}
	for len(lines) > 0 {
		isDirectiveLine, err := parseFileDirective(directive, lines[0])
		if err != nil {
			return directive, lines, err
		}
		if !isDirectiveLine {
			return directive, lines, nil
		}
		lines = lines[1:]
	}
	return directive, lines, nil
}

// https://www.onlinetool.io/gitoembed/widget?url=https%3A%2F%2Fgithub.com%2Fessentialbooks%2Fbooks%2Fblob%2Fmaster%2Fbooks%2Fgo%2F0020-basic-types%2Fbooleans.go
// to:
// books/go/0020-basic-types/booleans.go
// returns empty string if doesn't conform to what we expect
func gitoembedToRelativePath(uri string) string {
	parsed, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	switch parsed.Host {
	case "www.onlinetool.io", "onlinetool.io":
		// do nothing
	default:
		return ""
	}
	path := parsed.Path
	if path != "/gitoembed/widget" {
		return ""
	}
	uri = parsed.Query().Get("url")
	// https://github.com/essentialbooks/books/blob/master/books/go/0020-basic-types/booleans.go
	parsed, err = url.Parse(uri)
	if parsed.Host != "github.com" {
		return ""
	}
	path = strings.TrimPrefix(parsed.Path, "/essentialbooks/books/")
	if path == parsed.Path {
		return ""
	}
	// blob/master/books/go/0020-basic-types/booleans.go
	path = strings.TrimPrefix(path, "blob/")
	// master/books/go/0020-basic-types/booleans.go
	// those are branch names. Should I just strip first 2 elements from the path?
	path = strings.TrimPrefix(path, "master/")
	path = strings.TrimPrefix(path, "notion/")
	// books/go/0020-basic-types/booleans.go
	return path
}

// we don't want to show our // :show annotations in snippets
func removeAnnotationLines(lines []string) []string {
	var res []string
	prevWasEmpty := false
	for _, l := range lines {
		if strings.Contains(l, "// :show ") {
			continue
		}
		if len(l) == 0 && prevWasEmpty {
			continue
		}
		prevWasEmpty = len(l) == 0
		res = append(res, l)
	}
	return res
}

// convert local path like books/go/foo.go into path to the file in a github repo
func getGitHubPathForFile(path string) string {
	return "https://github.com/essentialbooks/books/blob/master/" + toUnixPath(path)
}

func setGoPlaygroundID(b *Book, sf *SourceFile) error {
	if sf.Lang != "go" {
		return nil
	}
	if sf.Directive.NoPlayground {
		return nil
	}
	id, err := getSha1ToGoPlaygroundIDCached(b, sf.DataToRun())
	if err != nil {
		return err
	}
	sf.GoPlaygroundID = id
	sf.PlaygroundURI = "https://goplay.space/#" + sf.GoPlaygroundID
	return nil
}

var allowedLanguages = map[string]bool{
	"go":         true,
	"javascript": true,
	"cpp":        true,
}

func setGlotPlaygroundID(b *Book, sf *SourceFile) error {
	if !sf.Directive.Glot {
		return nil
	}
	if sf.Directive.NoPlayground {
		return nil
	}

	lang := strings.ToLower(sf.Lang)
	lang = glotConvertLanguage(lang)
	if _, ok := allowedLanguages[lang]; !ok {
		return fmt.Errorf("'%s' ('%s') is not a supported language", sf.Lang, lang)
	}

	fileName := sf.Directive.FileName
	snippetName := sf.SnippetName
	id, err := getSha1ToGlotPlaygroundIDCached(b, sf.DataToRun(), snippetName, fileName, lang)
	if err != nil {
		return err
	}
	sf.GlotPlaygroundID = id
	sf.PlaygroundURI = "https://glot.io/snippets/" + sf.GlotPlaygroundID
	return nil
}

func setSourceFileData(sf *SourceFile, data []byte) error {
	sf.CodeSnippet = data
	sf.LinesRaw = dataToLines(sf.CodeSnippet)
	lines := sf.LinesRaw
	directive, lines, err := extractFileDirectives(lines)
	sf.Directive = directive
	sf.LinesToRun = removeAnnotationLines(lines)
	sf.LinesCode, err = extractCodeSnippets(lines)
	return err
}

func loadSourceFile(b *Book, path string) (*SourceFile, error) {
	data, err := common.ReadFileNormalized(path)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(path)
	lang := getLangFromFileExt(filepath.Ext(path))
	gitHubURL := getGitHubPathForFile(path)
	sf := &SourceFile{
		Path:      path,
		FileName:  name,
		Lang:      lang,
		GitHubURL: gitHubURL,
	}

	err = setSourceFileData(sf, data)
	if err != nil {
		fmt.Printf("loadSourceFile: '%s', setSourceFileData() failed with '%s'\n", path, err)
		panicIfErr(err)
	}
	if sf.Directive.NoOutput {
		fmt.Printf("NoOutput for '%s'\n", path)
	}
	setGoPlaygroundID(b, sf)
	err = getOutputCached(b, sf)
	fmt.Printf("loadSourceFile: '%s', lang: '%s'\n", path, lang)
	return sf, nil
}

// TODO: remove when all code moved to repl.it
func extractSourceFiles(b *Book, p *Page) {
	page := p.NotionPage
	for _, block := range page.Root.Content {
		if block.Type != notionapi.BlockEmbed {
			continue
		}
		uri := block.FormatEmbed.DisplaySource
		if strings.Contains(uri, "repl.it/") {
			panic("NYI")
			continue
		}
		relativePath := gitoembedToRelativePath(uri)
		if relativePath == "" {
			panic("NYI")
			fmt.Printf("Couldn't parse embed uri '%s'\n", uri)
			continue
		}
		panic("NYI")
		// fmt.Printf("Embed uri: %s, relativePath: %s\n", uri, relativePath)
		//path := filepath.Join(wd, relativePath)
		path := relativePath
		sf, err := loadSourceFile(b, path)
		if err != nil {
			fmt.Printf("extractSourceFiles: loadSourceFile('%s') (uri: '%s') failed with '%s'\n", path, uri, err)
			panicIfErr(err)
		}
		sf.EmbedURL = uri
		p.SourceFiles = append(p.SourceFiles, sf)
	}
}
