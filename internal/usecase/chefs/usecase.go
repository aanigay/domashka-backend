package chefs

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	entity "domashka-backend/internal/entity/chefs"
	dish "domashka-backend/internal/entity/dishes"
)

const (
	avatarFilePrefixTpl      = "avatar/m/chef_%d"
	smallAvatarFilePrefixTpl = "avatar/s/chef_%d"
)

type Usecase struct {
	chefRepo chefRepo
	geoRepo  geoRepo
	s3Client s3Client
}

func New(chefRepo chefRepo, repo geoRepo, s3client s3Client) *Usecase {
	return &Usecase{chefRepo: chefRepo, geoRepo: repo, s3Client: s3client}
}

func (u *Usecase) GetTopChefs(ctx context.Context, limit int) ([]entity.Chef, error) {
	return u.chefRepo.GetTopChefs(ctx, limit)
}

func (u *Usecase) GetNearestChefs(ctx context.Context, lat, long float64, distance, limit int) ([]entity.Chef, error) {
	chefs, err := u.chefRepo.GetNearestChefs(ctx, lat, long, distance, limit)
	if err != nil {
		return nil, err
	}
	for idx, chef := range chefs {
		c, err := u.chefRepo.GetChefRatingByChefID(ctx, chef.ID)
		if err != nil {
			return nil, err
		}
		chefs[idx].Rating = c.Rating
		chefs[idx].ReviewsCount = c.ReviewsCount
	}
	return chefs, nil
}

func (u *Usecase) GetDistanceToChef(ctx context.Context, lat, long float64, id int64) (float64, error) {
	return u.geoRepo.GetDistanceToChef(ctx, lat, long, id)
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
	filePrefix := fmt.Sprintf(avatarFilePrefixTpl, chefID)
	url, err := u.s3Client.UploadPicture(ctx, filePrefix, fileHeader)
	if err != nil {
		return "", err
	}
	err = u.chefRepo.SaveChefAvatar(ctx, chefID, url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (u *Usecase) UploadSmallAvatar(ctx context.Context, chefID int64, fileHeader *multipart.FileHeader) (string, error) {
	filePrefix := fmt.Sprintf(smallAvatarFilePrefixTpl, chefID)
	url, err := u.s3Client.UploadPicture(ctx, filePrefix, fileHeader)
	if err != nil {
		return "", err
	}
	err = u.chefRepo.SetSmallAvatar(ctx, chefID, url)
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

func (u *Usecase) GetAll(ctx context.Context) ([]entity.Chef, error) {
	chefs, err := u.chefRepo.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for idx, chef := range chefs {
		c, err := u.chefRepo.GetChefRatingByChefID(ctx, chef.ID)
		if err != nil {
			return nil, err
		}
		chefs[idx].Rating = c.Rating
		chefs[idx].ReviewsCount = c.ReviewsCount
	}
	return chefs, nil
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
