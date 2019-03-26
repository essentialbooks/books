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
	CodeFull  string
	CodeToRun string
	Lang      string

	GlotOutput string
	// id of glot.id code snippet for this code, if exists
	GlotID string

	// id of https://goplay.space snippet for this code, if exists
	GoPlayID string

	// sha1 of CodeFull, calculated on demand
	sha1 string
}

func (c *CodeInfo) Sha1() string {
	if c.sha1 == "" {
		c.sha1 = u.Sha1HexOfBytes([]byte(c.CodeFull))
	}
	return c.sha1
}

// a function in preparation for possible multiple
func (c *CodeInfo) Output() string {
	return c.GlotOutput
}

type Cache struct {
	// path of the cache file
	path string

	sha1ToCode map[string]*CodeInfo
}

func (c *Cache) getCodeInfoBySha1(sha1 string) (*CodeInfo, bool) {
	existing := true
	codeInfo := c.sha1ToCode[sha1]
	if codeInfo == nil {
		existing = false
		codeInfo = &CodeInfo{}
		c.sha1ToCode[sha1] = codeInfo
	}
	return codeInfo, existing
}

func (c *Cache) saveRecord(rec *siser.Record) error {
	f := openForAppend(c.path)
	defer f.Close()
	w := siser.NewWriter(f)
	_, err := w.WriteRecord(rec)
	return err
}

func (c *Cache) saveCodeInfo(code *CodeInfo) {
	panicIf(code.CodeFull == "")
	codeCurr, existing := c.getCodeInfoBySha1(code.Sha1())
	same := existing &&
		codeCurr.Lang == code.Lang &&
		codeCurr.CodeFull == code.CodeFull &&
		codeCurr.CodeToRun == code.CodeToRun

	if same {
		lg("saveCodeInfo: skipping because didn't change\n")
		return
	}
	rec := &siser.Record{
		Name: "code",
	}
	if code.Lang != "" {
		rec.Append("Lang", code.Lang)
		codeCurr.Lang = code.Lang
	}
	rec.Append("Sha1", code.Sha1())
	rec.Append("CodeFull", code.CodeFull)
	rec.Append("CodeToRun", code.CodeToRun)
	err := c.saveRecord(rec)
	must(err)

	codeCurr.Lang = code.Lang
	codeCurr.CodeFull = code.CodeFull
	codeCurr.CodeToRun = code.CodeToRun
}

func (c *Cache) saveGlotInfo(sha1, glotID, glotOutput string) {
	codeCurr, existing := c.getCodeInfoBySha1(sha1)
	same := existing &&
		codeCurr.GlotID == glotID &&
		codeCurr.GlotOutput == glotOutput
	if same {
		lg("saveGlotInfo: skipping because didn't change\n")
		return
	}

	rec := &siser.Record{
		Name: "glot",
	}
	panicIf(sha1 == "")
	rec.Append("Sha1", sha1)
	rec.Append("GlotID", glotID)
	rec.Append("GlotOutput", glotOutput)
	err := c.saveRecord(rec)
	must(err)
}

func (c *Cache) loadCodeInfo(rec *siser.Record) {
	sha1, ok := rec.Get("Sha1")
	panicIf(!ok)
	code, _ := c.getCodeInfoBySha1(sha1)

	panicIf(rec.Name != "code")
	for _, e := range rec.Entries {
		switch e.Key {
		case "CodeFull":
			panicIf(code.CodeFull != "" && code.CodeFull != e.Value)
			code.CodeFull = e.Value
		case "CodeToRun":
			code.CodeToRun = e.Value
		case "Lang":
			code.Lang = e.Value
		case "Sha1":
			// no-op
		default:
			must(fmt.Errorf("Unrecognized key '%s'", e.Key))
		}
	}
	panicIf(code.CodeFull == "")
	panicIf(sha1 != "" && sha1 != code.Sha1(), "sha1 != code.Sha1() (%s != %s)", sha1, code.Sha1())
}

func (c *Cache) loadGlot(rec *siser.Record) {
	sha1, ok := rec.Get("Sha1")
	panicIf(!ok)
	code, _ := c.getCodeInfoBySha1(sha1)

	panicIf(rec.Name != "code")
	for _, e := range rec.Entries {
		switch e.Key {
		case "GlotID":
			code.GlotID = e.Value
		case "GlotOutput":
			code.GlotOutput = e.Value
		case "Sha1":
			// no-op
		default:
			must(fmt.Errorf("Unrecognized key '%s'", e.Key))
		}
	}
	panicIf(code.CodeFull == "")
	panicIf(sha1 != "" && sha1 != code.Sha1(), "sha1 != code.Sha1() (%s != %s)", sha1, code.Sha1())
}

func loadCache(path string) *Cache {
	lg("loadCache: %s\n", path)
	dir := filepath.Dir(path)
	// the directory must exist
	_, err := os.Stat(dir)
	must(err)

	c := &Cache{
		path:       path,
		sha1ToCode: map[string]*CodeInfo{},
	}

	f, err := os.Open(path)
	if err != nil {
		// it's ok if file doesn't exist
		lg("  cache file %s doesn't exist\n", path)
		return c
	}
	defer f.Close()

	nRecords := 0
	r := siser.NewReader(bufio.NewReader(f))
	for r.ReadNextRecord() {
		rec := r.Record
		if rec.Name == "code" {
			c.loadCodeInfo(rec)
		} else if rec.Name == "glot" {
			c.loadGlot(rec)
		} else {
			panic(fmt.Errorf("unknown record type: '%s'", rec.Name))
		}
		nRecords++
	}
	must(r.Err())
	fmt.Printf(" got %d cache records\n", nRecords)
	return c
}
