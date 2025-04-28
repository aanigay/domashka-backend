package users

import (
	"context"

	usersEntity "domashka-backend/internal/entity/users"
)

type UseCase struct {
	repo repo
}

func New(repo repo) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Create(ctx context.Context, user *usersEntity.User) error {
	return u.repo.Create(ctx, user)
}

func (u *UseCase) GetByID(ctx context.Context, id int64) (*usersEntity.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) Update(ctx context.Context, id int64, user usersEntity.User) error {
	return u.repo.Update(ctx, id, user)
}

func (u *UseCase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
