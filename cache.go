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
	recNameGist       = "gist"       // content of the gist
	recNameGistOutput = "gistoutput" // result of evaluating the gist
)

type Cache struct {
	// path of the cache file
	path string

	// id of the gist to gist content. Might be outdated if
	// gist changes
	gistIDToGist         map[string]string
	gistSha1ToGistOutput map[string]string
}

func (c *Cache) saveRecord(rec *siser.Record) error {
	f := openForAppend(c.path)
	defer u.CloseNoError(f)
	w := siser.NewWriter(f)
	_, err := w.WriteRecord(rec)
	return err
}

func recGetMust(rec *siser.Record, name string) string {
	v, ok := rec.Get(name)
	u.PanicIf(!ok)
	return v
}

func recGetMustNonEmpty(rec *siser.Record, name string) string {
	v, ok := rec.Get(name)
	u.PanicIf(!ok || v == "")
	return v
}

func (c *Cache) saveGist(gistID, gist string) {
	rec := &siser.Record{
		Name: recNameGist,
	}
	u.PanicIf(gistID == "" || gist == "")
	rec.Append("GistID", gistID)
	rec.Append("Gist", gist)
	err := c.saveRecord(rec)
	must(err)
	c.gistIDToGist[gistID] = gist
}

func (c *Cache) loadGist(rec *siser.Record) {
	u.PanicIf(rec.Name != recNameGist)
	gistID := recGetMustNonEmpty(rec, "GistID")
	gist := recGetMustNonEmpty(rec, "Gist")
	c.gistIDToGist[gistID] = gist
}

func (c *Cache) saveGistOutput(gist, output string) {
	rec := &siser.Record{
		Name: recNameGistOutput,
	}
	//TODO: probably remove, it's ok to have no output
	//u.PanicIf(output == "")
	rec.Append("Gist", gist)
	rec.Append("GistOutput", output)
	err := c.saveRecord(rec)
	must(err)
	sha1 := u.Sha1HexOfBytes([]byte(gist))
	c.gistSha1ToGistOutput[sha1] = output
}

func (c *Cache) loadGistOutput(rec *siser.Record) {
	gist := recGetMustNonEmpty(rec, "Gist")
	output := recGetMust(rec, "GistOutput")
	sha1 := u.Sha1HexOfBytes([]byte(gist))
	c.gistSha1ToGistOutput[sha1] = output
}

func loadCache(path string) *Cache {
	logf("loadCache: %s\n", path)
	dir := filepath.Dir(path)
	// the directory must exist
	_, err := os.Stat(dir)
	must(err)

	c := &Cache{
		path:                 path,
		gistIDToGist:         map[string]string{},
		gistSha1ToGistOutput: map[string]string{},
	}

	f, err := os.Open(path)
	if err != nil {
		// it's ok if file doesn't exist
		logf("  cache file %s doesn't exist\n", path)
		return c
	}
	defer u.CloseNoError(f)

	nRecords := 0
	r := siser.NewReader(bufio.NewReader(f))
	for r.ReadNextRecord() {
		rec := r.Record
		switch rec.Name {
		case recNameGlotOutput:
			// not used anymore
		case recNameGlotID:
			// not used anymore
		case recNameGist:
			c.loadGist(rec)
		case recNameGistOutput:
			c.loadGistOutput(rec)
		default:
			panic(fmt.Errorf("unknown record type: '%s'", rec.Name))
		}
		nRecords++
	}
	must(r.Err())
	logf(" got %d cache records\n", nRecords)
	return c
}
