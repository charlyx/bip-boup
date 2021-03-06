package main

import (
	"io/ioutil"
	"log"
	"os"
)

// Loggers
var (
	logFlags = log.Ldate | log.Ltime | log.Lshortfile

	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLog(processName string, isDebug bool) {
	Debug = log.New(ioutil.Discard, "["+processName+"] DEBUG: ", logFlags)
	if isDebug {
		Debug = log.New(os.Stdout, "["+processName+"] DEBUG: ", logFlags)
	}
	Info = log.New(os.Stdout, "["+processName+"] INFO: ", logFlags)
	Warning = log.New(os.Stdout, "["+processName+"] WARNING: ", logFlags)
	Error = log.New(os.Stdout, "["+processName+"] ERROR: ", logFlags)
}
