package cjsqldriver

import (
	"database/sql/driver"
)

type CjSqlTx struct {
	mc *cjConnectionProxy
	driver.Tx
}

func (tx *CjSqlTx) Commit() (err error) {
	sqlDriverLogger.Debug("Commit tx")
	if err := tx.Tx.Commit(); err != nil {
		return err
	}

	tx.mc.inTransaction = false
	tx.mc.useSourceConn = nil
	return nil
}

func (tx *CjSqlTx) Rollback() (err error) {
	sqlDriverLogger.Debug("Rollback tx")
	if err := tx.Tx.Rollback(); err != nil {
		return err
	}

	tx.mc.inTransaction = false
	tx.mc.useSourceConn = nil
	return nil
}
