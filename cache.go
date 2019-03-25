package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kjk/siser"
)

type CodeSnippetWithOutput struct {
	// code as extracted from
	Code       string
	Lang       string
	GlotOutput string
	// id of glot.id code snippet for this code, if exists
	GlotID string
	// id of https://goplay.space snippet for this code, if exists
	GoPlayID string

	// calculated from Code
	codeSha1 string
}

func (c *CodeSnippetWithOutput) Output() string {
	// we might support output from other source in the future
	return c.GlotOutput
}

type Cache struct {
	path string

	sha1ToCode map[string]*CodeSnippetWithOutput
}

func (c *Cache) addCodeSnippet(snippet *CodeSnippetWithOutput) error {
	// TODO: write me:
	// - calc sha1 of Code
	// - if doesn't exist or different than current, save to a file
	//   changed should only happen if we expand what we do
	//fmt.Printf("addCodeSnippet: %s => %s\n", sha1, id)
	/*
		r := siser.Record{
			Keys:   []string{"sha1", "id"},
			Values: []string{sha1, id},
			Name:   "goplayid",
		}
		f := openForAppend(c.path)
		defer f.Close()
		w := siser.NewWriter(f)
		_, err := w.WriteRecord(&r)
		return err
	*/
	return nil
}

func loadCache(path string) *Cache {
	lg("loadCache: %s\n", path)
	dir := filepath.Dir(path)
	// the directory must exist
	_, err := os.Stat(dir)
	must(err)

	cache := &Cache{
		path:       path,
		sha1ToCode: map[string]*CodeSnippetWithOutput{},
	}

	f, err := os.Open(path)
	if err != nil {
		// it's ok if file doesn't exist
		lg("Cache file %s doesn't exist\n", path)
		return cache
	}
	defer f.Close()

	r := siser.NewReader(bufio.NewReader(f))
	for r.ReadNextRecord() {
		rec := r.Record
		if rec.Name == "code" {
			panic("NYI")
		} else {
			panic(fmt.Errorf("unknown record: '%s'", rec.Name))
		}
	}
	must(r.Err())
	return cache
}
