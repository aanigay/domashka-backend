package reviews

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"

	reviewEntity "domashka-backend/internal/entity/reviews"
)

// Usecase инкапсулирует бизнес-логику работы с отзывами.
type Usecase struct {
	repo                reviewRepo
	userRepo            userRepo
	orderRepo           orderRepo
	dishReviewsProducer KafkaWriter
	chefReviewsProducer KafkaWriter
}

var ErrReviewNotFound = errors.New("review not found")

// New создает новый экземпляр бизнес-слоя для отзывов.
func New(repo reviewRepo, userRepo userRepo, ordersRepo orderRepo, dishReviewsProducer, chefReviewsProducer KafkaWriter) *Usecase {
	return &Usecase{repo: repo, userRepo: userRepo, dishReviewsProducer: dishReviewsProducer, chefReviewsProducer: chefReviewsProducer, orderRepo: ordersRepo}
}

func (u *Usecase) CreateReview(ctx context.Context, review reviewEntity.Review) error {
	// переиспользуем у себя Create(ctx, *Review)
	return u.Create(ctx, &review)
}

func (u *Usecase) GetReviewsByChefID(ctx context.Context, chefID int64) ([]reviewEntity.Review, error) {
	// includeOnly=true, без лимита
	return u.ListByChef(ctx, chefID, true, nil)
}

// GetByID возвращает отзыв по его идентификатору.
func (u *Usecase) GetByID(ctx context.Context, id int64) (*reviewEntity.Review, error) {
	rv, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rv.IsDeleted {
		return nil, ErrReviewNotFound
	}
	return rv, nil
}

// ListByChef возвращает список отзывов для данного шефа.
// includeOnly = true — только отзывы, помеченные include_in_rating.
// limit — максимальное число отзывов (nil = без ограничения).
func (u *Usecase) ListByChef(ctx context.Context, chefID int64, includeOnly bool, limit *int) ([]reviewEntity.Review, error) {
	list, err := u.repo.ListByChef(ctx, chefID, includeOnly, limit)
	if err != nil {
		return nil, err
	}
	// Фильтруем мягко удалённые
	var out []reviewEntity.Review
	for _, rv := range list {
		if !rv.IsDeleted {
			out = append(out, rv)
		}
	}
	return out, nil
}

// Create создает новый отзыв, проверяя корректность полей.
func (u *Usecase) Create(ctx context.Context, rv *reviewEntity.Review) error {
	// Валидируем звезды
	if rv == nil {
		return fmt.Errorf("review is nil")
	}
	if rv.Stars < 1 || rv.Stars > 5 {
		return errors.New("stars must be between 1 and 5")
	}
	// todo Сменить IsVerified на false
	rv.IsVerified = true
	rv.IncludeInRating = false
	rv.IsDeleted = false

	err := u.repo.Create(ctx, rv)
	if err != nil {
		return err
	}
	type chefReview struct {
		ChefID int64 `json:"chef_id"`
		Rating int16 `json:"rating"`
	}
	type dishReview struct {
		DishID   int64 `json:"dish_id"`
		ReviewID int64 `json:"review_id"`
		Rating   int16 `json:"rating"`
	}
	r := &chefReview{
		ChefID: rv.ChefID,
		Rating: rv.Stars,
	}

	msg, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = u.chefReviewsProducer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(rv.ID, 10)),
		Value: msg,
	})
	if err != nil {
		return err
	}

	items, err := u.orderRepo.GetCartItemsByOrderID(ctx, rv.OrderID)
	if err != nil {
		return err
	}
	dishReviews := map[int64]dishReview{}
	for _, item := range items {
		dishReviews[item.Dish.ID] = dishReview{
			DishID:   item.Dish.ID,
			ReviewID: rv.ID,
			Rating:   rv.Stars,
		}
	}
	dishRatingMessages := make([]kafka.Message, 0, len(dishReviews))
	for dishID, review := range dishReviews {
		msg, err := json.Marshal(review)
		if err != nil {
			return err
		}
		dishRatingMessages = append(dishRatingMessages, kafka.Message{
			Key:   []byte(fmt.Sprintf("%s_%s", strconv.FormatInt(rv.ID, 10), strconv.FormatInt(dishID, 10))),
			Value: msg,
		})
	}
	err = u.dishReviewsProducer.WriteMessages(ctx, dishRatingMessages...)
	if err != nil {
		return err
	}
	return nil
}

// Update изменяет отзыв. Изменяются только контент и флаги.
func (u *Usecase) Update(ctx context.Context, rv *reviewEntity.Review) error {
	// Убедимся, что отзыв существует и не удалён
	if rv == nil {
		return fmt.Errorf("review is nil")
	}
	existing, err := u.repo.GetByID(ctx, rv.ID)
	if err != nil {
		return err
	}
	if existing.IsDeleted {
		return ErrReviewNotFound
	}
	// Валидируем звезды
	if rv.Stars < 1 || rv.Stars > 5 {
		return errors.New("stars must be between 1 and 5")
	}
	// При обновлении не меняем created_at; updated_at проставится в репозитории
	return u.repo.Update(ctx, rv)
}

// SoftDelete помечает отзыв как удалённый.
func (u *Usecase) SoftDelete(ctx context.Context, id int64) error {
	// Проверим, что отзыв существует
	rv, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rv.IsDeleted {
		return ErrReviewNotFound
	}
	return u.repo.SoftDelete(ctx, id)
}
func (u *Usecase) GetReviewByOrderAndUserID(ctx context.Context, chefID, userID int64) (*reviewEntity.Review, error) {
	return u.repo.GetReviewByOrderAndUserID(ctx, chefID, userID)
}

// GetFullReviewsByChefID возвращает сразу ReviewWithUser
func (u *Usecase) GetFullReviewsByChefID(ctx context.Context, chefID int64, limit int) ([]reviewEntity.ReviewWithUser, error) {
	reviews, err := u.repo.ListByChef(ctx, chefID, false, &limit)
	if err != nil {
		return nil, err
	}

	// 2) для каждого отзыва достаём юзера и собираем ReviewWithUser
	out := make([]reviewEntity.ReviewWithUser, 0, len(reviews))
	for _, rv := range reviews {
		user, err := u.userRepo.GetByID(ctx, rv.UserID)
		if err != nil {
			continue
		}
		out = append(out, reviewEntity.ReviewWithUser{
			Review:   rv,
			UserName: user.FirstName,
		})
	}
	return out, nil
}

func (u *Usecase) GetRatingStats(ctx context.Context, chefID int64) (map[int16]int32, error) {
	reviews, err := u.ListByChef(ctx, chefID, false, nil)
	if err != nil {
		return nil, err
	}

	ratingMap := make(map[int16]int32)
	for _, review := range reviews {
		ratingMap[review.Stars]++
	}
	return ratingMap, nil
}
