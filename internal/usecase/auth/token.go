package auth

type TokenGenerator interface {
	Generate(userID int64) (string, error)
}
