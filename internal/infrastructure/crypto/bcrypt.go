package crypto

import (
	"golang.org/x/crypto/bcrypt"

	domainuser "github.com/Nurdiansyah15/ddd-arch/internal/domain/master/user"
)

// compile-time assertion
var _ domainuser.PasswordHasher = (*BcryptHasher)(nil)

// BcryptHasher implements domain.PasswordHasher menggunakan bcrypt.
type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func (h *BcryptHasher) Check(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
