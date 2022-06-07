package cjsqldriver

import (
	"code.com/tars/goframework/kissgo/appzaplog"
	"code.com/tars/goframework/kissgo/appzaplog/zap"
	"code.com/tars/goframework/kissgo/appzaplog/zap/zapcore"
	_ "code.com/tars/goframework/tars/servant"
)

var sqlDriverLogger = new(cjSqlDriverLogger)

// sql driver logger
type cjSqlDriverLogger struct {
}

func (s *cjSqlDriverLogger) Debug(msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("business", "cjsqldriver"))
	appzaplog.Debug(msg, fields...)
}

func (s *cjSqlDriverLogger) Info(msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("business", "cjsqldriver"))
	appzaplog.Info(msg, fields...)
}

func (s *cjSqlDriverLogger) Error(msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("business", "cjsqldriver"))
	appzaplog.Error(msg, fields...)
}

func (s *cjSqlDriverLogger) Warn(msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("business", "cjsqldriver"))
	appzaplog.Warn(msg, fields...)
}
