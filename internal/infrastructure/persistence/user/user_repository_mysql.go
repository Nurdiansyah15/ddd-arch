package user

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryMySQL struct {
	DB *sqlx.DB
}

func (r *UserRepositoryMySQL) FindByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.DB.Get(&u, `
		SELECT id, email, password_hash, is_active
		FROM users
		WHERE email = ?
	`, email)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
