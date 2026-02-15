package user

import (
	userdomain "github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
)

type ListUsecase struct {
	Repo userdomain.Repository
}

func NewListUsecase(repo userdomain.Repository) *ListUsecase {
	return &ListUsecase{Repo: repo}
}

type ListResponseItem struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (uc *ListUsecase) Execute() ([]*ListResponseItem, error) {
	users, err := uc.Repo.List()
	if err != nil {
		return nil, err
	}
	out := make([]*ListResponseItem, 0, len(users))
	for _, u := range users {
		out = append(out, &ListResponseItem{ID: u.ID, Email: u.Email})
	}
	return out, nil
}
