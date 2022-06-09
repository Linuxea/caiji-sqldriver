package cjsqldriver

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-sql-driver/mysql"
)

func init() {
	cjSqlDriver := &CJSqlDriver{Driver: &mysql.MySQLDriver{}}
	sql.Register("cjmysql", cjSqlDriver)
}

// dsn config
type Dsn struct {
	ReadWeight  int    `json:"read_weight"`  // 读权重
	WriteWeight int    `json:"write_weight"` // 写权重
	Dsn         string `json:"dsn"`          // standard dsn
	Flag        string `json:"flag"`         // flag
}

var dbResolverPolicy Policy = &weightPolicy{R: rand.New(rand.NewSource(time.Now().UnixNano()))}

func SetDbResolverPolicy(policy Policy) {
	dbResolverPolicy = policy
}

type CJSqlDriver struct {
	driver.Driver
}

func (d CJSqlDriver) Open(dsn string) (driver.Conn, error) {

	var dsns []*Dsn
	err := json.Unmarshal([]byte(dsn), &dsns)
	if err != nil {
		dsns = make([]*Dsn, 0, 1)
		// 构造主从连接
		dsns = append(dsns, &Dsn{
			ReadWeight:  1,
			WriteWeight: 1,
			Dsn:         dsn,
		})
	}

	total := make([]*connectionItem, 0, len(dsns))
	write := make([]*connectionItem, 0)

	for index := range dsns {

		tmp := dsns[index]
		conn, err := d.Driver.Open(tmp.Dsn)
		if err != nil {
			return nil, err
		}

		newConn := &connectionItem{
			Conn:        conn,
			readWeight:  tmp.ReadWeight,
			writeWeight: tmp.WriteWeight,
			flag:        tmp.Flag,
		}

		total = append(total, newConn)

		if tmp.WriteWeight > 0 {
			write = append(write, newConn)
		}

	}

	return &cjConnectionProxy{all: total, write: write, policy: dbResolverPolicy}, nil
}
