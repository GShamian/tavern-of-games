package store

import "github.com/GShamian/tavern-of-games/internal/app/model"

// UserRepository interface
type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}
