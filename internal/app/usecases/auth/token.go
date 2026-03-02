package auth

// TokenGenerator provides token generation operations used by usecases.
type TokenGenerator interface {
	GenerateAccess(userID int64) (string, error)
	GenerateRefresh(userID int64) (string, error)
}
