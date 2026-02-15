package user

import "errors"

var (
	ErrUserInactive           = errors.New("user is inactive")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	IsActive     bool
}

func (u *User) Authenticate(password string) error {
	if !u.IsActive {
		return ErrUserInactive
	}

	if !CheckPassword(password, u.PasswordHash) {
		return ErrInvalidPassword
	}

	return nil
}
