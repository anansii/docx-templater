// see http://www.goinggo.net/2013/11/using-log-package-in-go.html for ideas

package main

import (
	"log"
	"os"
)

var (
	Msg   *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func setupLogging() {
	Msg = log.New(os.Stdout, "", 0)
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
