package cjsqldriver

import (
	"database/sql/driver"
)

type CjSqlTx struct {
	mc *cjConnectionProxy
	driver.Tx
}

func (tx *CjSqlTx) Commit() (err error) {
	CjSqlDriverLogger.Log("commit", tx.mc.useSourceConn.flag)
	if err := tx.Tx.Commit(); err != nil {
		return err
	}

	tx.mc.inTransaction = false
	tx.mc.useSourceConn = nil
	return nil
}

func (tx *CjSqlTx) Rollback() (err error) {
	CjSqlDriverLogger.Log("rollback", tx.mc.useSourceConn.flag)
	if err := tx.Tx.Rollback(); err != nil {
		return err
	}

	tx.mc.inTransaction = false
	tx.mc.useSourceConn = nil
	return nil
}
