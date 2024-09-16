package domain

import (
	"github.com/mathehluiz/plant-care-tracker/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id         int64  `json:"-"`
	ExternalId string `json:"external_id"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`

	Active   bool `json:"active"`
	Verified bool `json:"verified"`

	Roles []string `json:"roles"`
}

func NewUser(username, email, password string, roles []string) (*User, error) {
	m := &User{
		Email:    email,
		Username: username,
		Active:   true,
		Roles:    roles,
	}

	if len(m.Username) < 4 || len(m.Username) > 20 {
		return nil, errs.ErrInvalidUsername
	}

	if len(m.Email) < 4 || len(m.Email) > 100 {
		return nil, errs.ErrInvalidEmail
	}

	if err := m.HashPass(password); err != nil {
		return nil, err
	}

	return m, nil
}

func (u *User) HashPass(str string) error {
	if len(str) < 8 || len(str) > 32 {
		return errs.ErrInvalidPassword
	}

	u.Password = str
	return u.hashPassword()
}

func (u *User) hashPassword() error {
	p, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(p)
	return nil
}

func (u *User) VerifyPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
	if err != nil {
		return errs.ErrInvalidPassword
	}
	return nil
}
