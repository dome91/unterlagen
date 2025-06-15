package administration

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type UserMessages interface {
	PublishUserCreated(user User) error
	SubscribeUserCreated(subscriber func(user User) error) error
}

type User struct {
	Username                string
	Password                string
	Role                    UserRole
	PasswordChangeNecessary bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (user User) IsValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}

type UserRepository interface {
	Save(user User) error
	FindByUsername(username string) (User, error)
	FindAll() ([]User, error)
	FindAllByRole(role UserRole) ([]User, error)
	ExistsByRole(role UserRole) bool
}

type users struct {
	repository UserRepository
	messages   UserMessages
}

func (users *users) CreateUser(username string, password string, role UserRole) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Username:                username,
		Password:                string(hashedPassword),
		Role:                    role,
		PasswordChangeNecessary: role == UserRoleUser,
	}

	err = users.repository.Save(user)
	if err != nil {
		return err
	}

	return users.messages.PublishUserCreated(user)
}

func (users *users) AdminExists() bool {
	return users.repository.ExistsByRole(UserRoleAdmin)
}

func (users *users) GetUser(username string) (User, error) {
	return users.repository.FindByUsername(username)
}

func (users *users) GetAllUsers() ([]User, error) {
	return users.repository.FindAll()
}

func (users *users) GetAllUsersByRole(role UserRole) ([]User, error) {
	return users.repository.FindAllByRole(role)
}

func newUsers(repository UserRepository, messages UserMessages) *users {
	return &users{
		repository: repository,
		messages:   messages,
	}
}
