package sqlite

import (
	"database/sql"
	"errors"
	"unterlagen/features/administration"

	"github.com/jmoiron/sqlx"
)

var _ administration.UserRepository = &UserRepository{}

type UserRepository struct {
	db *sqlx.DB
}

// ExistsByRole implements administration.UserRepository.
func (u *UserRepository) ExistsByRole(role administration.UserRole) bool {
	var count int
	err := u.db.Get(&count, "SELECT COUNT(*) FROM users WHERE role = ?", string(role))
	if err != nil {
		return false
	}
	return count > 0
}

// FindAll implements administration.UserRepository.
func (u *UserRepository) FindAll() ([]administration.User, error) {
	var users []administration.User
	err := u.db.Select(&users, "SELECT * FROM users order by username")
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindAllByRole implements administration.UserRepository.
func (u *UserRepository) FindAllByRole(role administration.UserRole) ([]administration.User, error) {
	var users []administration.User
	err := u.db.Select(&users, "SELECT * FROM users WHERE role = ?", string(role))
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindByUsername implements administration.UserRepository.
func (u *UserRepository) FindByUsername(username string) (administration.User, error) {
	var user administration.User
	err := u.db.Get(&user, "SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return administration.User{}, errors.New("user not found")
		}
		return administration.User{}, err
	}
	return user, nil
}

// Save implements administration.UserRepository.
func (u *UserRepository) Save(user administration.User) error {
	query := `
		INSERT INTO users (username, password, role, password_change_necessary, created_at, updated_at)
		VALUES (:username, :password, :role, :password_change_necessary, datetime(), datetime())
		ON CONFLICT(username) DO UPDATE SET
		password = excluded.password,
		role = excluded.role,
		password_change_necessary = excluded.password_change_necessary,
		updated_at = datetime()
	`
	_, err := u.db.NamedExec(query, user)
	return err
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}
