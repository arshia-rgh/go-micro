package data

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
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

func (u *PostgresTestRepository) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `SELECT * FROM users order by last_name`

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User

		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, err

}

func (u *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `SELECT * FROM users WHERE email = ?`

	row := db.QueryRowContext(ctx, query, email)
	var user User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, err

}

func (u *PostgresTestRepository) GetById(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `SELECT * FROM users WHERE id = ?`

	row := db.QueryRowContext(ctx, query, id)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (u *PostgresTestRepository) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `UPDATE users SET 
                 email = ?,
                 first_name = ?,
                 last_name = ?,
                 user_active = ?,
                 updated_at = ?
             where id = ?
             `

	_, err := db.ExecContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)

	return err
}

// Delete delete by User.ID
func (u *PostgresTestRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `DELETE FROM users WHERE id = ?`

	_, err := db.ExecContext(ctx, query, id)

	return err
}

func (u *PostgresTestRepository) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (email, first_name, last_name, password, user_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`

	var newID int
	row := db.QueryRowContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPass,
		user.Active,
		time.Now(),
		time.Now(),
	)

	err = row.Scan(&newID)

	return newID, err
}

func (u *PostgresTestRepository) ChangePassword(password string, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = ? where id = ?`

	_, err = db.ExecContext(ctx, query, hashedPass, user.ID)

	return err

}

func (u *PostgresTestRepository) PasswordsMatches(plainPassword string, user User) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))

	if err != nil {
		switch {

		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil

		default:
			return false, err
		}
	}

	return true, nil
}
