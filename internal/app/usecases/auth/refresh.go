package auth

import (
	"errors"
)

var ErrInvalidRefresh = errors.New("invalid refresh token")

type RefreshUsecase struct {
	TokenSvc TokenService
}

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (uc *RefreshUsecase) Execute(req RefreshRequest) (*RefreshResponse, error) {
	uid, err := uc.TokenSvc.ValidateRefresh(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefresh
	}

	access, err := uc.TokenSvc.GenerateAccess(uid)
	if err != nil {
		return nil, err
	}

	return &RefreshResponse{AccessToken: access}, nil
}

// TokenService defines the minimal token operations required by the usecase.
type TokenService interface {
	ValidateRefresh(string) (int64, error)
	GenerateAccess(int64) (string, error)
}

func NewRefreshUsecase(ts TokenService) *RefreshUsecase {
	return &RefreshUsecase{TokenSvc: ts}
}
