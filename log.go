package main

import (
	"fmt"
	"os"
)

var (
	logFile *os.File
)

func fmtSmart(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

func openLog() func() {
	var err error
	logFile, err = os.Create("log.txt")
	must(err)
	return func() {
		_ = logFile.Close()
		logFile = nil
	}
}

func logf(format string, args ...interface{}) {
	s := fmtSmart(format, args...)
	fmt.Print(s)
	if logFile != nil {
		_, _ = fmt.Fprint(logFile, s)
	}
}

func logVerbose(format string, args ...interface{}) {
	if logFile == nil {
		return
	}
	s := fmtSmart(format, args...)
	_, _ = fmt.Fprintf(logFile, s)
}

func logFatalf(format string, args ...interface{}) {
	s := fmtSmart(format, args...)
	fmt.Print(s)
	if logFile != nil {
		_, _ = fmt.Fprint(logFile, s)
	}
	os.Exit(1)
}
