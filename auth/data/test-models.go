package data

import (
	"database/sql"
	"time"
)

type PostgresTestRepository struct {
	Conn *sql.DB
}

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{
		Conn: db,
	}
}

// GetAll returns a slice of all users, sorted by last name
func (u *PostgresTestRepository) GetAll() ([]*User, error) {
	users := []*User{}

	return users, nil
}

// GetByEmail returns one user by email
func (u *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		FirstName: "First",
		LastName:  "Last",
		Email:     "me@here.com",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// GetOne returns one user by id
func (u *PostgresTestRepository) GetById(id int) (*User, error) {
	user := User{
		ID:        1,
		FirstName: "First",
		LastName:  "Last",
		Email:     "me@here.com",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *PostgresTestRepository) Update(user User) error {
	return nil
}

func (u *PostgresTestRepository) Delete(id int) error {
	return nil
}

func (u *PostgresTestRepository) Insert(user User) (int, error) {
	return 2, nil
}

func (u *PostgresTestRepository) ChangePassword(password string, user User) error {
	return nil
}

func (u *PostgresTestRepository) PasswordsMatches(plainPassword string, user User) (bool, error) {
	return true, nil
}
