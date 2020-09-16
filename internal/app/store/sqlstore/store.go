package sqlstore

import (
	"database/sql"

	"github.com/GShamian/tavern-of-games/internal/app/store"
	// import postgres
	_ "github.com/lib/pq"
)

// Store object, that is made to store information about DB
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New func. Constructor for Store object
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
