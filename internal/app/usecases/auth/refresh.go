package auth

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
)

// TokenService defines the minimal token operations required by the usecase.
type TokenService interface {
	ValidateRefresh(string) (int64, error)
	GenerateAccess(int64) (string, error)
}

type RefreshUsecase struct {
	TokenSvc TokenService
}

func NewRefreshUsecase(ts TokenService) *RefreshUsecase {
	return &RefreshUsecase{TokenSvc: ts}
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
		return nil, apperror.Unauthorized("invalid refresh token")
	}

	access, err := uc.TokenSvc.GenerateAccess(uid)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &RefreshResponse{AccessToken: access}, nil
}
