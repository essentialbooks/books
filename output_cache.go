package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/kjk/u"
)

// runs `go run ${path}` and returns captured output`
func getGoOutput(path string) (string, error) {
	dir, fileName := filepath.Split(path)
	cmd := exec.Command("go", "run", fileName)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func getRunCmdOutput(path string, runCmd string) (string, error) {
	parts, err := shlex.Split(runCmd)
	maybePanicIfErr(err)
	if err != nil {
		return "", err
	}
	exeName := parts[0]
	parts = parts[1:]
	var parts2 []string
	srcDir, srcFileName := filepath.Split(path)

	// remove empty lines and replace variables
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		switch part {
		case "$file":
			part = srcFileName
		}
		parts2 = append(parts2, part)
	}
	//fmt.Printf("getRunCmdOutput: running '%s' with args '%#v'\n", exeName, parts2)
	cmd := exec.Command(exeName, parts2...)
	cmd.Dir = srcDir
	out, err := cmd.CombinedOutput()
	//fmt.Printf("getRunCmdOutput: out:\n%s\n", string(out))
	return string(out), err
}

func stripCurrentPathFromOutput(s string) string {
	path, err := filepath.Abs(".")
	u.PanicIfErr(err)
	return strings.Replace(s, path, "", -1)
}

// it executes a code file and captures the output
// optional runCmd says
func getOutput(path string, runCmd string) (string, error) {
	if runCmd != "" {
		//fmt.Printf("Found :run cmd '%s' in '%s'\n", runCmd, path)
		s, err := getRunCmdOutput(path, runCmd)
		return stripCurrentPathFromOutput(s), err
	}

	// do default
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".go" {
		s, err := getGoOutput(path)
		return stripCurrentPathFromOutput(s), err
	}
	return "", fmt.Errorf("getOutput(%s): files with extension '%s' are not supported", path, ext)
}

// for a given file, get output of executing this command
// We cache this as it is the most expensive part of rebuilding books
// If allowError is true, we silence an error from executed command
// This is useful when e.g. executing "go run" on a program that is
// intentionally not valid.
func getOutputCached(b *Book, sf *SourceFile) error {
	// TODO: make the check for when not to execute the file
	// even better
	if sf.Directive.NoOutput {
		if sf.Directive.Glot {
			// if running on glot, we want to execute even if
			// we don't show the output (to check syntax errors)
		} else {
			return nil
		}
	}
	code := sf.DataToRun()
	sha1Hex := u.Sha1HexOfBytes(code)

	cof := b.cache.sha1ToCode[sha1Hex]
	if cof != nil {
		// is guaranteed to exist
		sf.Output = cof.Output()
		return nil
	}
	if flgNoUpdateOutput {
		return nil
	}

	if sf.Directive.Glot {
		f := &glotFile{
			Name:    sf.Directive.FileName,
			Content: string(code),
		}
		if sf.Directive.RunCmd != "" {
			fmt.Printf("  run command: %s\n", sf.Directive.RunCmd)
		}
		req := &glotRunRequest{
			Command:  sf.Directive.RunCmd,
			language: sf.Lang,
			Files:    []*glotFile{f},
		}
		rsp, err := glotRun(req)
		panicIfErr(err)
		s := rsp.Stdout + rsp.Stderr
		if rsp.Error != "" {
			if !sf.Directive.AllowError {
				//fmt.Printf("getOutput('%s'), output is:\n%s\n", path, s)
				return errors.New(rsp.Stderr)
			}
		} else if rsp.Stderr != "" {
			fmt.Printf("getOutputCached: got stderr: %s\n", rsp.Stderr)
			if !sf.Directive.AllowError {
				//fmt.Printf("getOutput('%s'), output is:\n%s\n", path, s)
				return errors.New(rsp.Stderr)
			}
		}
		sf.Output = s
		fmt.Printf("Got glot output (%d bytes) for %s from %s\n", len(sf.Output), sf.SnippetName, sf.NotionOriginURL)
	} else {
		path := sf.Path
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".json", ".csv", ".yml", ".xml":
			return nil
		}

		// fmt.Printf("loadFileCached('%s') failed with '%s'\n", outputPath, err)
		s, err := getOutput(path, sf.Directive.RunCmd)
		if err != nil {
			if !sf.Directive.AllowError {
				fmt.Printf("getOutput('%s'), got error:\n%s\n", path, s)
				return err
			}
			err = nil
		}
		sf.Output = s
		fmt.Printf("Got output (%d bytes) for '%s' by running locally\n", len(sf.Output), path)
	}

	panic("NYI")
	/*
		cof = getCurrentOutputCacheFile(b)
		cof.doc = kvstore.ReplaceOrAppend(cof.doc, sha1Hex, sf.Output)

		b.sha1ToCachedOutputFile[sha1Hex] = cof
		saveCachedOutputFiles(b)
	*/
	return nil
}
