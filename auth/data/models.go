package data

import "database/sql"

var db *sql.DB

type Models struct {
	User User
}

type User struct {
}

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{User: User{}}

}
