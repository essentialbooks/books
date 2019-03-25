package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kjk/siser"
	"github.com/kjk/u"
)

type CodeInfo struct {
	// code as extracted from
	CodeFull   string
	CodeToRun  string
	Lang       string
	GlotOutput string
	// id of glot.id code snippet for this code, if exists
	GlotID string
	// id of https://goplay.space snippet for this code, if exists
	GoPlayID string

	sha1 string
}

func (c *CodeInfo) Sha1() string {
	if c.sha1 == "" {
		c.sha1 = u.Sha1HexOfBytes([]byte(c.CodeFull))
	}
	return c.sha1
}

type Cache struct {
	path string

	sha1ToCode map[string]*CodeInfo
}

func codeInfoToRecord(code *CodeInfo) *siser.Record {
	rec := &siser.Record{
		Name: "code",
	}
	rec.Append("CodeFull", code.CodeFull)
	rec.Append("CodeToRun", code.CodeToRun)
	if code.Lang != "" {
		rec.Append("Lang", code.Lang)
	}
	if code.GlotID != "" {
		rec.Append("GlotID", code.GlotID)
	}
	if code.GoPlayID != "" {
		rec.Append("GoPlayID", code.GoPlayID)
	}
	return rec
}

func recordToCodeInfo(rec *siser.Record) *CodeInfo {
	code := &CodeInfo{}
	panicIf(rec.Name != "code")
	for _, e := range rec.Entries {
		switch e.Key {
		case "CodeFull":
			code.CodeFull = e.Value
		case "CodeToRun":
			code.CodeToRun = e.Value
		case "Lang":
			code.Lang = e.Value
		case "GlotID":
			code.GlotID = e.Value
		case "GoPlayID":
			code.GoPlayID = e.Value
		default:
			err := fmt.Errorf("Unrecognized key '%s'", e.Key)
			must(err)
		}
	}
	panicIf(code.CodeFull == "")
	return code
}

func (c *Cache) saveCodeInfo(code *CodeInfo) error {
	c.sha1ToCode[code.Sha1()] = code
	rec := codeInfoToRecord(code)
	//lg("addCodeSnippet from https://notion.so/%s\n", id)
	f := openForAppend(c.path)
	defer f.Close()
	w := siser.NewWriter(f)
	_, err := w.WriteRecord(rec)
	return err
}

func loadCache(path string) *Cache {
	lg("loadCache: %s", path)
	dir := filepath.Dir(path)
	// the directory must exist
	_, err := os.Stat(dir)
	must(err)

	cache := &Cache{
		path:       path,
		sha1ToCode: map[string]*CodeInfo{},
	}

	f, err := os.Open(path)
	if err != nil {
		// it's ok if file doesn't exist
		lg("Cache file %s doesn't exist\n", path)
		return cache
	}
	defer f.Close()

	nRecords := 0
	r := siser.NewReader(bufio.NewReader(f))
	for r.ReadNextRecord() {
		rec := r.Record
		if rec.Name == "code" {
			codeInfo := recordToCodeInfo(rec)
			cache.sha1ToCode[codeInfo.Sha1()] = codeInfo
		} else {
			panic(fmt.Errorf("unknown record type: '%s'", rec.Name))
		}
		nRecords++
	}
	must(r.Err())
	fmt.Printf(" got %d cache records\n", nRecords)
	return cache
}
