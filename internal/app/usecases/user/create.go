package user

import "github.com/Nurdiansyah15/ddd-arch/internal/app/domain/master/user"

type CreateUsecase struct {
	Repo        user.Repository
	UserService *user.UserService // domain service
}

func NewCreateUsecase(repo user.Repository, svc *user.UserService) *CreateUsecase {
	return &CreateUsecase{Repo: repo, UserService: svc}
}

type CreateRequest struct {
	Email    string
	Password string
}

type CreateResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (uc *CreateUsecase) Execute(req CreateRequest) (*CreateResponse, error) {
	// Domain Service: cek apakah email sudah dipakai
	if err := uc.UserService.CheckEmailAvailability(req.Email); err != nil {
		return nil, err
	}

	hash, err := user.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Email:        req.Email,
		PasswordHash: hash,
		IsActive:     true,
	}

	if err := uc.Repo.Create(u); err != nil {
		return nil, err
	}

	return &CreateResponse{ID: u.ID, Email: u.Email}, nil
}
