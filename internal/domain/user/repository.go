package user

type Repository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id int64) (*User, error)
	Create(u *User) error
	Update(u *User) error
	Delete(id int64) error
	List() ([]*User, error)
}
