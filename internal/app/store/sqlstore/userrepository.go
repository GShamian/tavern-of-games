package sqlstore

import (
	"database/sql"

	"github.com/GShamian/tavern-of-games/internal/app/model"
	"github.com/GShamian/tavern-of-games/internal/app/store"
)

// UserRepository object for storing store entities
type UserRepository struct {
	store *Store
}

// Create func. Writing an email and encrypted password in the fields
// in DB that match to imported User
func (r *UserRepository) Create(u *model.User) error {
	// Checking user's fields for incorrect entries
	if err := u.Validate(); err != nil {
		return err
	}
	// Creating encrypted password. Chech user.go documentation
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	// Writing an email and ecrypted password in DB
	return r.store.db.QueryRow("INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

// FindByEmail func. Finding user with the right (email we need) email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}
