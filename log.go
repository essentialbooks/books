package main

import (
	"fmt"
	"os"
)

var (
	logFile *os.File
)

func openLog() {
	var err error
	logFile, err = os.Create("log.txt")
	must(err)
}

func closeLog() {
	if logFile == nil {
		return
	}
	logFile.Close()
	logFile = nil
}

func lg(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	if logFile != nil {
		fmt.Fprint(logFile, s)
	}
	fmt.Print(s)
}

func verbose(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	if logFile != nil {
		fmt.Fprint(logFile, s)
	}
}
