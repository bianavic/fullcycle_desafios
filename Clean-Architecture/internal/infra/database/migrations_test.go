package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

func TestMigrations(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/orders")
	assert.NoError(t, err)
	defer db.Close()

	var tableName string
	err = db.QueryRow("SHOW TABLES LIKE 'orders'").Scan(&tableName)
	assert.NoError(t, err)
	assert.Equal(t, "orders", tableName)

	rows, err := db.Query("DESCRIBE orders")
	assert.NoError(t, err)
	defer rows.Close()

	columns := make(map[string]string)
	for rows.Next() {
		var field, typ, null, key, extra string
		var def *string
		err = rows.Scan(&field, &typ, &null, &key, &def, &extra)
		assert.NoError(t, err)
		columns[field] = typ
	}

	assert.Equal(t, "varchar(255)", columns["id"])
	assert.Equal(t, "float", columns["price"])
	assert.Equal(t, "float", columns["tax"])
	assert.Equal(t, "float", columns["final_price"])
}
