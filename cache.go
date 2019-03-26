package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kjk/siser"
	"github.com/kjk/u"
)

const (
	recNameGlotOutput = "glotoutput"
	recNameGlotID     = "glotid"
)

type Cache struct {
	// path of the cache file
	path string

	sha1ToGlotOutput map[string]*GlotOutput
	// id of glot.id code snippet for this code, if exists
	sha1ToGlotID map[string]string
	// id of https://goplay.space snippet for this code, if exists
	sha1ToGoPlayID map[string]string
}

func NewCache(path string) *Cache {
	return &Cache{
		path:             path,
		sha1ToGlotOutput: map[string]*GlotOutput{},
		sha1ToGlotID:     map[string]string{},
		sha1ToGoPlayID:   map[string]string{},
	}
}

func (c *Cache) saveRecord(rec *siser.Record) error {
	f := openForAppend(c.path)
	defer f.Close()
	w := siser.NewWriter(f)
	_, err := w.WriteRecord(rec)
	return err
}

type GlotOutput struct {
	Lang      string
	FileName  string
	CodeFull  string
	CodeToRun string
	RunCmd    string
	Output    string

	// sha1 of CodeFull, calculated on demand
	sha1 string
}

func (c *GlotOutput) Sha1() string {
	if c.sha1 == "" {
		c.sha1 = u.Sha1HexOfBytes([]byte(c.CodeFull))
	}
	return c.sha1
}

func (c *Cache) saveGlotOutput(code *GlotOutput) {
	panicIf(code.CodeFull == "")
	panicIf(c.sha1ToGlotOutput[code.Sha1()] != nil)
	rec := &siser.Record{
		Name: recNameGlotOutput,
	}
	rec.Append("Sha1", code.Sha1())
	rec.Append("Lang", code.Lang)
	rec.Append("FileName", code.FileName)
	rec.Append("CodeFull", code.CodeFull)
	rec.Append("CodeToRun", code.CodeToRun)
	if code.RunCmd != "" {
		rec.Append("RunCmd", code.RunCmd)
	}
	rec.Append("Output", code.Output)
	err := c.saveRecord(rec)
	must(err)
	c.sha1ToGlotOutput[code.Sha1()] = code
}

func (c *Cache) loadGlotOutput(rec *siser.Record) {
	panicIf(rec.Name != recNameGlotOutput)
	sha1, ok := rec.Get("Sha1")
	panicIf(!ok || sha1 == "")
	panicIf(c.sha1ToGlotOutput[sha1] != nil)

	o := &GlotOutput{}
	o.Lang, ok = rec.Get("Lang")
	panicIf(!ok)
	o.FileName, ok = rec.Get("FileName")
	panicIf(!ok)
	o.CodeFull, ok = rec.Get("CodeFull")
	panicIf(!ok)
	o.CodeToRun, ok = rec.Get("CodeToRun")
	panicIf(!ok)
	o.RunCmd, ok = rec.Get("RunCmd")
	o.Output, ok = rec.Get("Output")
	panicIf(!ok)

	panicIf(o.CodeFull == "")
	panicIf(sha1 != o.Sha1(), "sha1 != code.Sha1() (%s != %s)", sha1, o.Sha1())
	c.sha1ToGlotOutput[sha1] = o
}

func (c *Cache) saveGlotID(sha1, glotID string) {
	panicIf(c.sha1ToGlotID[sha1] != "")
	rec := &siser.Record{
		Name: recNameGlotID,
	}
	panicIf(sha1 == "")
	rec.Append("Sha1", sha1)
	rec.Append("GlotID", glotID)
	err := c.saveRecord(rec)
	must(err)
	c.sha1ToGlotID[sha1] = glotID
}

func (c *Cache) loadGlotID(rec *siser.Record) {
	panicIf(rec.Name != recNameGlotID)
	sha1, ok := rec.Get("Sha1")
	panicIf(!ok || sha1 == "")
	panicIf(c.sha1ToGlotID[sha1] != "")

	glotid, ok := rec.Get("GlotID")
	panicIf(!ok)
	c.sha1ToGlotID[sha1] = glotid
}

func loadCache(path string) *Cache {
	lg("loadCache: %s\n", path)
	dir := filepath.Dir(path)
	// the directory must exist
	_, err := os.Stat(dir)
	must(err)

	c := NewCache(path)

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
		if rec.Name == recNameGlotOutput {
			c.loadGlotOutput(rec)
		} else if rec.Name == recNameGlotID {
			c.loadGlotID(rec)
		} else {
			panic(fmt.Errorf("unknown record type: '%s'", rec.Name))
		}
		nRecords++
	}
	must(r.Err())
	fmt.Printf(" got %d cache records\n", nRecords)
	return c
}
