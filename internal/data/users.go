package data

import (
	"errors"
	"time"

	"github.com/arian-nj/site/back/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser

}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid email address")
}

func ValidatePasswordPlainText(v *validator.Validator, plainPassword string) {
	v.Check(plainPassword != "", "password", "must be provided")
	v.Check(len(plainPassword) >= 8, "password", "must be at least 8 character long")
	v.Check(len(plainPassword) <= 72, "password", "must not be more than 72 character long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500")
	v.Check(len(user.Name) >= 5, "name", "must be more than 5 character")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlainText(v, *user.Password.Plaintext)
	}
	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}

}

type Password struct {
	Plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
