package main

// GetGlotPlaygroundID gets go playground id from content
func GetGlotPlaygroundID(b *Book, d []byte, snippetName string, fileName string, lang string) (string, bool, error) {
	panic("NYI")
	return "", false, nil
	/*
		sha1 := u.Sha1HexOfBytes(d)
		id, ok := b.cache.sha1ToGlotID[sha1]
		if ok {
			//fmt.Printf("GetGlotPlaygroundID: got %s from cache\n", sha1)
			return id, true, nil
		}

		//fmt.Printf("GetGlotPlaygroundID: no %s in cache of %d entries\n", sha1, len(b.cache.sha1ToGlotID))
		rsp, err := glotGetSnippedID(d, snippetName, fileName, lang)
		if err != nil {
			return "", false, err
		}
		err = b.cache.addGlotSha1ToID(sha1, rsp.ID)
		if err != nil {
			return "", false, err
		}
		return id, false, nil
	*/
}

func getSha1ToGlotPlaygroundIDCached(b *Book, d []byte, snippetName string, fileName string, lang string) (string, error) {
	panic("NYI")
	return "", nil
	/*
		id, fromCache, err := GetGlotPlaygroundID(b, d, snippetName, fileName, lang)
		if err != nil {
			return "", err
		}
		if !fromCache {
			sha1 := u.Sha1HexOfBytes(d)
			uri := "https://glot.io/snippets/" + id
			fmt.Printf("getSha1ToGlotPlaygroundIDCached: generated glot id %s => %s, %s\n", sha1, id, uri)
		}
		return id, nil
	*/
}
