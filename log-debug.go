// +build debug

package main

import (
	"log"
	"os"
)

var (
	Trace *log.Logger
)

func setupDebugLogging() {
	Trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
}
