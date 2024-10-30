package data

type Repository interface {
	GetAll() ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetById(id int) (*User, error)
	Update(user User) error
	Delete(id int) error
	Insert(user User) (int, error)
	ChangePassword(password string, user User) error
	PasswordsMatches(plainPassword string, user User) (bool, error)
}
