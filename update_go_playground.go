package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/kjk/u"

	"github.com/essentialbooks/books/pkg/common"
)

// submit the data to Go playground and get share id
func getGoPlaygroundShareID(d []byte) (string, error) {
	uri := "https://play.golang.org/share"
	r := bytes.NewBuffer(d)
	resp, err := http.Post(uri, "text/plain", r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http.Post returned error code '%s'", err)
	}
	d, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(d)), nil
}

func testGetGoPlaygroundShareIDAndExit() {
	path := "books/go/0230-mutex/rwlock.go"
	d, err := common.ReadFileNormalized(path)
	panicIfErr(err)
	shareID, err := getGoPlaygroundShareID(d)
	panicIfErr(err)
	fmt.Printf("share id: '%s'\n", shareID)
	os.Exit(0)
}

// GetGoPlaygroundID gets go playground id from content
func GetGoPlaygroundID(c *Cache, d string) (string, bool, error) {
	panic("NYI")
	return "", false, nil
	/*
		sha1 := u.Sha1HexOfBytes(d)
		id, ok := c.sha1ToGoPlayID[sha1]
		if ok {
			return id, true, nil
		}
		id, err := getGoPlaygroundShareID(d)
		if err != nil {
			return "", false, err
		}
		err = c.addGoPlaySha1ToID(sha1, id)
		if err != nil {
			return "", false, err
		}
		return id, false, nil
	*/
}

func getSha1ToGoPlaygroundIDCached(b *Book, d string) (string, error) {
	id, fromCache, err := GetGoPlaygroundID(b.cache, d)
	if err == nil && !fromCache {
		sha1 := u.Sha1HexOfBytes([]byte(d))
		fmt.Printf("getSha1ToGoPlaygroundIDCached: %s => %s\n", sha1, id)
	}
	return id, err
}
