package user

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

type userRepositoryPG struct {
	DB *sqlx.DB
}

func NewUserRepositoryPG(db *sqlx.DB) user.Repository {
	return &userRepositoryPG{DB: db}
}

func (r *userRepositoryPG) FindByEmail(email string) (*user.User, error) {
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

func (r *userRepositoryPG) Create(user *user.User) error {
	return r.DB.QueryRowx(`
		INSERT INTO users (email, password_hash, is_active)
		VALUES ($1, $2, $3)
		RETURNING id
	`, user.Email, user.PasswordHash, user.IsActive).
		Scan(&user.ID)
}

func (r *userRepositoryPG) FindByID(id int64) (*user.User, error) {
	var u user.User
	err := r.DB.Get(&u, `
		SELECT id, email, password_hash, is_active
		FROM users
		WHERE id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepositoryPG) Update(u *user.User) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET email = $1, password_hash = $2, is_active = $3
		WHERE id = $4
	`, u.Email, u.PasswordHash, u.IsActive, u.ID)
	return err
}

func (r *userRepositoryPG) Delete(id int64) error {
	_, err := r.DB.Exec(`
		DELETE FROM users WHERE id = $1
	`, id)
	return err
}

func (r *userRepositoryPG) List() ([]*user.User, error) {
	var users []*user.User
	err := r.DB.Select(&users, `
		SELECT id, email, password_hash, is_active
		FROM users
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	return users, nil
}
