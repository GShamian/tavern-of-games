package teststore

import (
	"github.com/GShamian/tavern-of-games/internal/app/model"
	"github.com/GShamian/tavern-of-games/internal/app/store"
)

// Store object for tests only
type Store struct {
	userRepository *UserRepository
}

// New func. Empty constructor (default constructor) for testing
// only store entities
func New() *Store {
	return &Store{}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.userRepository
}
