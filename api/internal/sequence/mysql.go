package sequence

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const sqlReplaceStub = `REPLACE INTO SEQUENCE(STUB) VALUES ('a')`

type Mysql struct {
	conn sqlx.SqlConn
}

func NewMysql(dsn string) Sequence {
	return &Mysql{
		conn: sqlx.NewMysql(dsn),
	}
}

func (m *Mysql) Next() (seq uint64, err error) {
	stmt, err := m.conn.Prepare(sqlReplaceStub)
	if err != nil {
		logx.Errorf("prepare stub sql failed,err:%v", err)
		return 0, err
	}
	defer stmt.Close()
	rest, err := stmt.Exec()
	if err != nil {
		return 0, err
	}
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		return
	}
	return uint64(lid), nil
}
