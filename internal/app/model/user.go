package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

// User object that has id, email, password and encrypted password fields
type User struct {
	ID                int
	Email             string
	Password          string
	EncryptedPassword string
}

// Validate func. Validating user instance for id, email and password
func (u *User) Validate() error {

	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

// BeforeCreate func. Encrypting password func that encrypts password and writes encrypted
// version in User's EncryptedPassword field
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc

	}
	return nil
}

// encryptString func. Using function GenerateFromPassword from bcrypt encrypts imported
// string (original unencrypted password) to encrypted variation and
// returns it
func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
