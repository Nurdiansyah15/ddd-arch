package user

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

type userRepositoryMySQL struct {
	DB *sqlx.DB
}

func NewUserRepositoryMySQL(db *sqlx.DB) user.Repository {
	return &userRepositoryMySQL{DB: db}
}

func (r *userRepositoryMySQL) FindByEmail(email string) (*user.User, error) {
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

func (r *userRepositoryMySQL) Create(user *user.User) error {
	res, err := r.DB.Exec(`
		INSERT INTO users (email, password_hash, is_active)
		VALUES (?, ?, ?)
	`, user.Email, user.PasswordHash, user.IsActive)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *userRepositoryMySQL) FindByID(id int64) (*user.User, error) {
	var u user.User
	err := r.DB.Get(&u, `
		SELECT id, email, password_hash, is_active
		FROM users
		WHERE id = ?
	`, id)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepositoryMySQL) Update(u *user.User) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET email = ?, password_hash = ?, is_active = ?
		WHERE id = ?
	`, u.Email, u.PasswordHash, u.IsActive, u.ID)
	return err
}

func (r *userRepositoryMySQL) Delete(id int64) error {
	_, err := r.DB.Exec(`
		DELETE FROM users WHERE id = ?
	`, id)
	return err
}

func (r *userRepositoryMySQL) List() ([]*user.User, error) {
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
