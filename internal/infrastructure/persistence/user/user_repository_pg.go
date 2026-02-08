package user

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryPG struct {
	DB *sqlx.DB
}

func (r *UserRepositoryPG) FindByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.DB.Get(&u, `
		SELECT id, email, password_hash, is_active
		FROM users
		WHERE email = $1
	`, email)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryPG) Create(user *user.User) error {
	return r.DB.QueryRowx(`
		INSERT INTO users (email, password_hash, is_active)
		VALUES ($1, $2, $3)
		RETURNING id
	`, user.Email, user.PasswordHash, user.IsActive).
		Scan(&user.ID)
}
