package chefs

import (
	"context"
	"log"
	"mime/multipart"

	entity "domashka-backend/internal/entity/chefs"
	dish "domashka-backend/internal/entity/dishes"
)

type Usecase struct {
	chefRepo chefRepo
}

func New(chefRepo chefRepo) *Usecase {
	return &Usecase{chefRepo: chefRepo}
}

func (u *Usecase) GetTopChefs(ctx context.Context, limit int) ([]entity.Chef, error) {
	return u.chefRepo.GetTopChefs(ctx, limit)
}

func (u *Usecase) GetChefByDishID(ctx context.Context, dishID int64) (*entity.Chef, error) {
	chef, err := u.chefRepo.GetChefByDishID(ctx, dishID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	chefRating, err := u.chefRepo.GetChefRatingByChefID(ctx, chef.ID)
	if err != nil {
		return nil, err
	}
	chef.Rating = chefRating.Rating
	chef.ReviewsCount = chefRating.ReviewsCount
	return chef, nil
}
func (u *Usecase) UploadAvatar(ctx context.Context, chefID int64, fileHeader *multipart.FileHeader) (string, error) {
	url, err := u.chefRepo.SaveChefAvatar(ctx, chefID, fileHeader)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (u *Usecase) GetChefByID(ctx context.Context, chefID int64) (*entity.Chef, error) {
	chef, err := u.chefRepo.GetChefByID(ctx, chefID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	chefRating, err := u.chefRepo.GetChefRatingByChefID(ctx, chef.ID)
	if err != nil {
		return nil, err
	}
	chef.Rating = chefRating.Rating
	chef.ReviewsCount = chefRating.ReviewsCount
	return chef, nil
}

func (u *Usecase) GetChefAvatarURLByDishID(ctx context.Context, dishID int64) (string, error) {
	return u.chefRepo.GetChefAvatarURLByDishID(ctx, dishID)
}

func (u *Usecase) GetChefAvatarURLByChefID(ctx context.Context, chefID int64) (string, error) {
	return u.chefRepo.GetChefAvatarURLByChefID(ctx, chefID)
}
func (u *Usecase) GetChefExperienceYears(ctx context.Context, chefID int64) (int, error) {
	return u.chefRepo.GetChefExperienceYears(ctx, chefID)
}
func (u *Usecase) GetChefCertifications(ctx context.Context, chefID int64) ([]entity.Certification, error) {
	return u.chefRepo.GetChefCertifications(ctx, chefID)
}
func (u *Usecase) GetDishesByChefID(ctx context.Context, chefID int64) ([]dish.Dish, error) {
	return u.chefRepo.GetDishesByChefID(ctx, chefID)
}
