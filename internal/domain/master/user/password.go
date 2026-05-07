package user

// PasswordHasher adalah port (interface) domain untuk operasi hashing password.
// Implementasinya ada di layer infrastructure (e.g. bcrypt).
type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
}
