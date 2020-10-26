package infrastructure

import (
	"fmt"
	"os"
	"testing"
)

func Test_SqliteHandler(t *testing.T) {
	h := NewSqliteHandler("test.db")
	if err := h.Execute("DROP TABLE IF EXISTS foo"); err != nil {
		t.Error(err)
	}
	if err := h.Execute("CREATE TABLE foo (id integer, name varchar(42))"); err != nil {
		t.Error(err)
	}
	if err := h.Execute("INSERT INTO foo (id, name) VALUES (23, 'johndoe')"); err != nil {
		t.Error(err)
	}
	row := h.Query("SELECT id, name FROM foo LIMIT 1")
	var id int
	var name string
	row.Next()
	if err := row.Scan(&id, &name); err != nil {
		t.Error(err)
	}
	if id != 23 {
		fmt.Println(id)
		t.Error()
	}
	if err := os.Remove("test.db"); err != nil {
		t.Error(err)
	}
}
