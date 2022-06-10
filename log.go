package cjsqldriver

import (
	"log"
	"os"
)

type Logger interface {
	Print(v ...interface{})
}

var CjSqlDriverLogger = NewcjSqlDriverLogger()

var defaultLogger Logger = log.New(os.Stdout, "\r\n", 0)

func NewcjSqlDriverLogger() *cjSqlDriverLogger {
	return &cjSqlDriverLogger{
		Logger: defaultLogger,
	}
}

type cjSqlDriverLogger struct {
	Logger
}

func (c *cjSqlDriverLogger) Log(v ...interface{}) {
	c.Print(v)
}

func (c *cjSqlDriverLogger) SetLogger(l Logger) {
	c.Logger = l
}
