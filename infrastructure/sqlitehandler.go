package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteHandler struct {
	Conn *sql.DB
}

func (handler *SqliteHandler) Execute(statement string) (err error) {
	_, err = handler.Conn.Exec(statement)
	return err
}

func (handler *SqliteHandler) Query(statement string) interfaces.Row {
	//fmt.Println(statement)
	rows, err := handler.Conn.Query(statement)
	if err != nil {
		fmt.Println(err)
		return new(SqliteRow)
	}
	row := new(SqliteRow)
	row.Rows = rows
	return row
}

type SqliteRow struct {
	Rows *sql.Rows
}

func (r SqliteRow) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r SqliteRow) Next() bool {
	return r.Rows.Next()
}

func NewSqliteHandler(dbFilename string) *SqliteHandler {
	conn, _ := sql.Open("sqlite3", dbFilename)
	sqliteHandler := new(SqliteHandler)
	sqliteHandler.Conn = conn
	return sqliteHandler
}
