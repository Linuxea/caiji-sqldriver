package cjsqldriver

import (
	"github.com/jinzhu/gorm"
)

func init() {
	dialect, ok := gorm.GetDialect("mysql")
	if !ok {
		panic("获取不到 mysql dialect")
	}

	gorm.RegisterDialect("cjmysql", dialect)
}
