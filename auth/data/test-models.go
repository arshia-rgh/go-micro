package data

import "database/sql"

type PostgresTestRepository struct {
	Conn *sql.DB
}
