package cjsqldriver

import (
	"log"
	"os"
)

type Logger interface {
	Print(v ...interface{})
}

var defaultLogger Logger = log.New(os.Stdout, "\r\n", 0)
