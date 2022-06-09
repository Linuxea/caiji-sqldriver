package cjsqldriver

import (
	"database/sql/driver"
	"strings"
)

type connectionItem struct {
	driver.Conn
	readWeight  int
	writeWeight int
	flag        string
}

type cjConnectionProxy struct {
	all           []*connectionItem
	write         []*connectionItem
	policy        Policy
	inTransaction bool
	useSourceConn *connectionItem // 保存事务中使用的连接
}

func (c *cjConnectionProxy) Prepare(query string) (driver.Stmt, error) {

	defaultLogger.Print("prepare sql", query)

	var useReadConn *connectionItem
	// 不在事务中
	if !c.inTransaction {
		if rawSQL := strings.TrimSpace(query); len(rawSQL) > 10 && strings.EqualFold(rawSQL[:6], "select") && !strings.EqualFold(rawSQL[len(rawSQL)-10:], "for update") {
			// read
			if len(c.all) == 1 {
				useReadConn = c.all[0]
			} else {
				useReadConn = c.policy.ResolveRead(c.all)
			}

			return useReadConn.Prepare(query)
		}
	}

	// 在事务中, 找到事务中使用的连接
	if c.useSourceConn != nil {
		return c.useSourceConn.Prepare(query)
	}

	// 找不到使用的连接 可能是特殊 sql 如:show index from xxx 不需要开启事务等等 或者 bug
	var useWriteConn *connectionItem
	if len(c.write) == 1 {
		useWriteConn = c.write[0]
	} else {
		useWriteConn = c.policy.ResolveWrite(c.write)
	}

	return useWriteConn.Prepare(query)
}

func (c *cjConnectionProxy) Close() error {

	defaultLogger.Print("Close")
	for index := range c.all {
		if err := c.all[index].Close(); err != nil {
			return err
		}
	}

	return nil
}

func (c *cjConnectionProxy) Begin() (driver.Tx, error) {
	defaultLogger.Print("Begin")
	conn := c.policy.ResolveWrite(c.write)
	tx, err := conn.Begin()
	if err != nil {
		return tx, err
	}

	// 标记为事务中
	c.inTransaction = true
	// 保存使用的事务连接
	c.useSourceConn = conn
	// wrapper tx
	myTx := &CjSqlTx{
		Tx: tx,
		mc: c,
	}

	return myTx, err
}
