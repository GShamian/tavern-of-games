package teststore

import (
	"github.com/GShamian/tavern-of-games/internal/app/store"

	"github.com/GShamian/tavern-of-games/internal/app/model"
)

// UserRepository object for testing only
type UserRepository struct {
	store *Store
	users map[int]*model.User
}

// Create func. Writing an email and encrypted password in the fields
// in DB that match to imported User. For additional information
// check userrepository.go documentation in sqlstore dir.
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

// FindByEmail func. Finding user with the right (email we need) email.
// Function for testing only purposes.
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// Find func. Finding user with the right (id we need) id.
// Function for testing only purposes.
func (r *UserRepository) Find(id int) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
