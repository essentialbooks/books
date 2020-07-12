package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kjk/u"
)

var must = u.Must
var panicIf = u.PanicIf

func logIfError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

// whitelisted characters valid in url
func validateRune(c rune) byte {
	if c >= 'a' && c <= 'z' {
		return byte(c)
	}
	if c >= '0' && c <= '9' {
		return byte(c)
	}
	if c == '-' || c == '_' || c == '.' || c == ' ' {
		return '-'
	}
	return 0
}

func charCanRepeat(c byte) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

// urlify generates safe url from tile by removing hazardous characters
func urlify(title string) string {
	s := strings.TrimSpace(title)
	s = strings.ToLower(s)
	var res []byte
	for _, r := range s {
		c := validateRune(r)
		if c == 0 {
			continue
		}
		// eliminute duplicate consequitive characters
		var prev byte
		if len(res) > 0 {
			prev = res[len(res)-1]
		}
		if c == prev && !charCanRepeat(c) {
			continue
		}
		res = append(res, c)
	}
	s = string(res)
	if len(s) > 128 {
		s = s[:128]
	}
	s = strings.TrimLeft(s, "-")
	s = strings.TrimRight(s, "-")
	return s
}

var (
	softErrorMode bool
	delayedErrors []string

	totalHTMLBytes         int
	totalHTMLBytesMinified int
)

func maybePanicIfErr(err error) {
	if err == nil {
		return
	}
	if !softErrorMode {
		u.Must(err)
	}
	delayedErrors = append(delayedErrors, err.Error())
}

func clearErrors() {
	delayedErrors = nil
	totalHTMLBytes = 0
	totalHTMLBytesMinified = 0
}

func printAndClearErrors() {
	fmt.Printf("HTML: optimized %d => %d (saved %d bytes)\n", totalHTMLBytes, totalHTMLBytesMinified, totalHTMLBytes-totalHTMLBytesMinified)
	if len(delayedErrors) == 0 {
		return
	}
	errStr := strings.Join(delayedErrors, "\n")
	fmt.Printf("\n%d errors:\n%s\n\n", len(delayedErrors), errStr)
	clearErrors()
}

func createDirForFileMaybeMust(path string) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	maybePanicIfErr(err)
}

func copyFileMaybeMust(dst, src string) error {
	createDirForFileMaybeMust(dst)
	err := u.CopyFile(dst, src)
	maybePanicIfErr(err)
	return err
}

// "foo.js" => "foo-${sha1}.js"
func nameToSha1Name(name, sha1Hex string) string {
	ext := filepath.Ext(name)
	n := len(name)
	s := name[:n-len(ext)]
	return s + "-" + sha1Hex[:8] + ext
}

func urlJoin(s1, s2 string) string {
	if strings.HasSuffix(s1, "/") {
		if strings.HasPrefix(s2, "/") {
			return s1 + s2[1:]
		}
		return s1 + s2
	}

	if strings.HasPrefix(s2, "/") {
		return s1 + s2
	}
	return s1 + "/" + s2
}

// removes empty lines from the beginning and end of the array
func trimEmptyLines(lines []string) []string {
	for len(lines) > 0 && len(lines[0]) == 0 {
		lines = lines[1:]
	}

	for len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	n := len(lines)
	res := make([]string, 0, n)
	prevWasEmpty := false
	for i := 0; i < n; i++ {
		l := lines[i]
		shouldAppend := l != "" || !prevWasEmpty
		prevWasEmpty = l == ""
		if shouldAppend {
			res = append(res, l)
		}
	}
	return res
}

func countStartChars(s string, c byte) int {
	for i := range s {
		if s[i] != c {
			return i
		}
	}
	return len(s)
}

// remove longest common space/tab prefix on non-empty lines
func shiftLines(lines []string) {
	maxTabPrefix := 1024
	maxSpacePrefix := 1024
	// first determine how much we can remove
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		n := countStartChars(line, ' ')
		if n > 0 {
			if n < maxSpacePrefix {
				maxSpacePrefix = n
			}
			continue
		}
		n = countStartChars(line, '\t')
		if n > 0 {
			if n < maxTabPrefix {
				maxTabPrefix = n
			}
			continue
		}
		// if doesn't start with space or tab, early abort
		return
	}
	if maxSpacePrefix == 1024 && maxTabPrefix == 1024 {
		return
	}

	toRemove := maxSpacePrefix
	if maxTabPrefix != 1024 {
		toRemove = maxTabPrefix
	}
	if toRemove == 0 {
		return
	}

	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		lines[i] = line[toRemove:]
	}
}

// replace potentially windows paths \foo\bar into unix paths /foo/bar
func toUnixPath(s string) string {
	return strings.Replace(s, `\`, "/", -1)
}

func dataToLines(d []byte) []string {
	s := string(d)
	return strings.Split(s, "\n")
}

func reverseStringSlice(a []string) {
	n := len(a) / 2
	for i := 0; i < n; i++ {
		a[i], a[n-i] = a[n-i], a[i]
	}
}

// turn "010 Defining a SetterGetter" to "Defining a SetterGetter"
func cleanTitle(s string) string {
	idx := strings.Index(s, " ")
	if idx == -1 {
		return s
	}
	if idx > 6 {
		return s
	}
	maybeNum := s[:idx]
	rest := s[idx+1:]
	_, err := strconv.Atoi(maybeNum)
	if err != nil {
		return s
	}
	return strings.TrimSpace(rest)
}

func openForAppend(path string) *os.File {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	must(err)
	return f
}

// for dev, rebuilds in background
func launchRollup(ctx context.Context) func() {
	path := filepath.Join("node_modules", ".bin", "rollup")
	if !u.FileExists(path) {
		// assume it's because node_modules doesn't exist, so install deps
		//cmd := exec.Command("npm", "install")
		cmd := exec.Command("yarn", "install")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		logf("Running: '%s'\n", cmd.String())
		err := cmd.Run()
		must(err)
	}
	cmd := exec.Command(path, "-c", "-w")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logf("Running: '%s'\n", cmd.String())
	err := cmd.Start()
	must(err)
	return func() {
		_ = cmd.Process.Kill()
		logf("Killed rollup\n")
	}
}

func buildFrontend() {
	{
		os.Remove("package-lock.json")
		os.RemoveAll("node_modules")
		cmd := exec.Command("yarn", "install")
		u.RunCmdMust(cmd)
	}
	// could also be
	// .\node_modules\.bin\rollup -c
	{
		cmd := exec.Command("yarn", "build-dev")
		u.RunCmdMust(cmd)
	}
}
