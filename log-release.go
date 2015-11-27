// +build !debug

package main

import (
	"io/ioutil"
	"log"
)

var (
	Trace *log.Logger
)

func setupDebugLogging() {
	Trace = log.New(ioutil.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
}
