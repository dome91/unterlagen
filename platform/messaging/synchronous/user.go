package synchronous

import (
	"log/slog"
	"unterlagen/features/administration"
)

var _ administration.UserMessages = &UserMessages{}

type UserMessages struct {
	userCreatedSubscribers []func(user administration.User) error
}

func (u *UserMessages) PublishUserCreated(user administration.User) error {
	for _, subscriber := range u.userCreatedSubscribers {
		err := subscriber(user)
		if err != nil {
			slog.Error("failed to process user created event", slog.String("error", err.Error()))
		}
	}
	return nil
}

func (u *UserMessages) SubscribeUserCreated(subscriber func(user administration.User) error) error {
	u.userCreatedSubscribers = append(u.userCreatedSubscribers, subscriber)
	return nil
}

func NewUserMessages() *UserMessages {
	return &UserMessages{
		userCreatedSubscribers: []func(user administration.User) error{},
	}
}
