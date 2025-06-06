package memory

import (
	"errors"
	"unterlagen/features/administration"
)

var _ administration.UserRepository = &UserRepository{}

type UserRepository struct {
	store map[string]administration.User
}

// ExistsByRole implements administration.UserRepository.
func (u *UserRepository) ExistsByRole(role administration.UserRole) bool {
	for _, user := range u.store {
		if user.Role == role {
			return true
		}
	}
	return false
}

// FindAllByRole implements administration.UserRepository.
func (u *UserRepository) FindAllByRole(role administration.UserRole) ([]administration.User, error) {
	var result []administration.User
	for _, user := range u.store {
		if user.Role == role {
			result = append(result, user)
		}
	}
	return result, nil
}

// FindByUsername implements administration.UserRepository.
func (u *UserRepository) FindByUsername(username string) (administration.User, error) {
	user, ok := u.store[username]
	if !ok {
		return administration.User{}, errors.New("user not found")
	}
	return user, nil
}

// Save implements administration.UserRepository.
func (u *UserRepository) Save(user administration.User) error {
	u.store[user.Username] = user
	return nil
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		store: make(map[string]administration.User),
	}
}
