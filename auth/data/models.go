package data

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

type Models struct {
	User User
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{User: User{}}

}

func (u *User) GetAll() ([]*User, error) {
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

func (u *User) GetByEmail(email string) (*User, error) {
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

func (u *User) GetById(id int) (*User, error) {
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

func (u *User) Update() error {
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
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)

	return err
}

// Delete delete by User.ID
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `DELETE FROM users WHERE id = ?`

	_, err := db.ExecContext(ctx, query, u.ID)

	return err
}

// DeleteById delete by the user id
func (u *User) DeleteById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `DELETE FROM users WHERE id = ?`

	_, err := db.ExecContext(ctx, query, id)

	return err
}

func (u *User) Insert() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (email, first_name, last_name, password, user_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`

	var newID int
	row := db.QueryRowContext(ctx, query,
		u.Email,
		u.FirstName,
		u.LastName,
		hashedPass,
		u.Active,
		time.Now(),
		time.Now(),
	)

	err = row.Scan(&newID)

	return newID, err
}

func (u *User) ChangePassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = ? where id = ?`

	_, err = db.ExecContext(ctx, query, hashedPass, u.ID)

	return err

}

func (u *User) PasswordsMatches(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))

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
